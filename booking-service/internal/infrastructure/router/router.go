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
	r := &Router{
		mux:            http.NewServeMux(),
		bookingHandler: bookingHandler,
		healthHandler:  healthHandler,
	}
	r.setupRoutes()
	return r
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) setupRoutes() {
	r.mux.Handle("/v1/", http.StripPrefix("/v1", r.mux))

	// Health check endpoint
	r.mux.HandleFunc("GET /health", r.withLogging(r.healthHandler.Check))

	// Booking endpoints
	r.mux.HandleFunc("POST /bookings", r.withLogging(r.bookingHandler.CreateBooking))
	r.mux.HandleFunc("GET /bookings", r.withLogging(r.bookingHandler.ListBookings))
	r.mux.HandleFunc("GET /bookings/{id}", r.withLogging(r.bookingHandler.GetBooking))
	r.mux.HandleFunc("DELETE /bookings/{id}", r.withLogging(r.bookingHandler.CancelBooking))
}

func (r *Router) withLogging(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next(w, req)
		duration := time.Since(start)
		// TODO: Replace with proper logging
		println(req.Method, req.URL.Path, duration)
	}
}
