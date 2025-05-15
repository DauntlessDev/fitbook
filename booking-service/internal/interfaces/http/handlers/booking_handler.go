package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/application/commands"
	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/application/queries"
)

type BookingHandler struct {
	createHandler *commands.CreateBookingHandler
	getHandler    *queries.GetBookingHandler
	listHandler   *queries.ListBookingsHandler
	cancelHandler *commands.CancelBookingHandler
}

func NewBookingHandler(
	createHandler *commands.CreateBookingHandler,
	getHandler *queries.GetBookingHandler,
	listHandler *queries.ListBookingsHandler,
	cancelHandler *commands.CancelBookingHandler,
) *BookingHandler {
	return &BookingHandler{
		createHandler: createHandler,
		getHandler:    getHandler,
		listHandler:   listHandler,
		cancelHandler: cancelHandler,
	}
}

func (handler *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreateBookingDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeBadRequest(w, "Invalid request body", err.Error())
		return
	}

	result, err := handler.createHandler.Handle(r.Context(), commands.CreateBookingCommand{DTO: &dto})
	if err != nil {
		handleBookingError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, result.Booking)
}

func (handler *BookingHandler) GetBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := r.PathValue("id")
	if bookingID == "" {
		writeBadRequest(w, "Booking ID is required")
		return
	}

	result, err := handler.getHandler.Handle(r.Context(), queries.GetBookingQuery{BookingID: bookingID})
	if err != nil {
		handleBookingError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result.Booking)
}

func (handler *BookingHandler) ListBookings(w http.ResponseWriter, r *http.Request) {
	var dto dtos.CreateBookingDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		writeBadRequest(w, "Invalid request body")
		return
	}

	if dto.UserID == "" && dto.GymID == "" {
		writeBadRequest(w, "Either user_id or gym_id must be provided")
		return
	}

	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		writeBadRequest(w, "Invalid start_time format")
		return
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		writeBadRequest(w, "Invalid end_time format")
		return
	}

	result, err := handler.listHandler.Handle(r.Context(), queries.ListBookingsQuery{
		UserID:    dto.UserID,
		GymID:     dto.GymID,
		StartTime: startTime,
		EndTime:   endTime,
	})

	if err != nil {
		writeInternalError(w)
		return
	}

	writeJSON(w, http.StatusOK, result.Bookings)
}

func (handler *BookingHandler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := r.PathValue("id")
	if bookingID == "" {
		writeBadRequest(w, "Booking ID is required")
		return
	}

	err := handler.cancelHandler.Handle(r.Context(), commands.CancelBookingCommand{BookingID: bookingID})
	if err != nil {
		handleBookingError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, nil)
}
