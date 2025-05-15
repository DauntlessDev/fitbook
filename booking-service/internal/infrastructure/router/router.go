package router

import (
	"net/http"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/interfaces/http/handlers"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type Router struct {
	mux            *http.ServeMux
	bookingHandler *handlers.BookingHandler
	healthHandler  *handlers.HealthHandler
}

func NewRouter(
	bookingHandler *handlers.BookingHandler,
	healthHandler *handlers.HealthHandler,
) *Router {
	router := &Router{
		mux:            http.NewServeMux(),
		bookingHandler: bookingHandler,
		healthHandler:  healthHandler,
	}
	router.setupRoutes()
	return router
}

func (router *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	router.mux.ServeHTTP(w, req)
}

func (router *Router) setupRoutes() {
	router.mux.Handle("/v1/", http.StripPrefix("/v1", router.mux))

	// Health check endpoint
	router.mux.HandleFunc("GET /health", router.withLogging(router.healthHandler.Check))

	// Booking endpoints
	router.mux.HandleFunc("POST /bookings", router.withLogging(router.bookingHandler.CreateBooking))
	router.mux.HandleFunc("GET /bookings", router.withLogging(router.bookingHandler.ListBookings))
	router.mux.HandleFunc("GET /bookings/{id}", router.withLogging(router.bookingHandler.GetBooking))
	router.mux.HandleFunc("DELETE /bookings/{id}", router.withLogging(router.bookingHandler.CancelBooking))
}

func (router *Router) withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next(w, req)
		duration := time.Since(start)
		// TODO: Replace with proper logging
		println(req.Method, req.URL.Path, duration)
	}
}
