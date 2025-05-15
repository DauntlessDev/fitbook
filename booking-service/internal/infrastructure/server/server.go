package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/yourusername/fitbook/booking-service/internal/application/commands"
	"github.com/yourusername/fitbook/booking-service/internal/application/queries"
	"github.com/yourusername/fitbook/booking-service/internal/infrastructure/config"
	"github.com/yourusername/fitbook/booking-service/internal/infrastructure/database"
	"github.com/yourusername/fitbook/booking-service/internal/infrastructure/events"
	"github.com/yourusername/fitbook/booking-service/internal/infrastructure/router"
	"github.com/yourusername/fitbook/booking-service/internal/interfaces/http/handlers"
)

func Start(cfg *config.Config) error {
	db, err := sql.Open("postgres", cfg.Database.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer db.Close()

	bookingRepo := database.NewBookingRepository(db)

	eventPublisher := events.NewEventPublisher()

	createBookingHandler := commands.NewCreateBookingHandler(bookingRepo, eventPublisher)
	cancelBookingHandler := commands.NewCancelBookingHandler(bookingRepo, eventPublisher)

	getBookingHandler := queries.NewGetBookingHandler(bookingRepo)
	listBookingsHandler := queries.NewListBookingsHandler(bookingRepo)

	bookingHandler := handlers.NewBookingHandler(
		createBookingHandler,
		getBookingHandler,
		listBookingsHandler,
		cancelBookingHandler,
	)
	healthHandler := handlers.NewHealthHandler()

	newRouter := router.NewRouter(bookingHandler, healthHandler)
	log.Println("Router initialized")

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port),
		Handler:      newRouter,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
	log.Printf("Server configured to listen on %s", srv.Addr)

	go func() {
		log.Printf("Starting server on %s...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return err
	}

	return nil
}
