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
	createHandler   *commands.CreateBookingHandler
	getHandler      *queries.GetBookingHandler
	listHandler     *queries.ListBookingsHandler
	cancelHandler   *commands.CancelBookingHandler
	confirmHandler  *commands.ConfirmBookingHandler
	completeHandler *commands.CompleteBookingHandler
}

func NewBookingHandler(
	createHandler *commands.CreateBookingHandler,
	getHandler *queries.GetBookingHandler,
	listHandler *queries.ListBookingsHandler,
	cancelHandler *commands.CancelBookingHandler,
	confirmHandler *commands.ConfirmBookingHandler,
	completeHandler *commands.CompleteBookingHandler,
) *BookingHandler {
	return &BookingHandler{
		createHandler:   createHandler,
		getHandler:      getHandler,
		listHandler:     listHandler,
		cancelHandler:   cancelHandler,
		confirmHandler:  confirmHandler,
		completeHandler: completeHandler,
	}
}

func (handler *BookingHandler) CreateBooking(writer http.ResponseWriter, request *http.Request) {
	var dto dtos.CreateBookingDTO
	if err := json.NewDecoder(request.Body).Decode(&dto); err != nil {
		writeBadRequest(writer, "Invalid request body", err.Error())
		return
	}

	result, err := handler.createHandler.Handle(request.Context(), commands.CreateBookingCommand{DTO: &dto})
	if err != nil {
		handleBookingError(writer, err)
		return
	}

	writeJSON(writer, http.StatusCreated, result.Booking)
}

func (handler *BookingHandler) GetBooking(writer http.ResponseWriter, request *http.Request) {
	bookingID := request.PathValue("id")
	if bookingID == "" {
		writeBadRequest(writer, "Booking ID is required")
		return
	}

	result, err := handler.getHandler.Handle(request.Context(), queries.GetBookingQuery{BookingID: bookingID})
	if err != nil {
		handleBookingError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, result.Booking)
}

func (handler *BookingHandler) ConfirmBooking(writer http.ResponseWriter, request *http.Request) {
	bookingID := request.PathValue("id")
	if bookingID == "" {
		writeBadRequest(writer, "Booking ID is required")
		return
	}

	err := handler.confirmHandler.Handle(request.Context(), commands.ConfirmBookingCommand{BookingID: bookingID})
	if err != nil {
		handleBookingError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, nil)
}

func (handler *BookingHandler) CompleteBooking(writer http.ResponseWriter, request *http.Request) {
	bookingID := request.PathValue("id")
	if bookingID == "" {
		writeBadRequest(writer, "Booking ID is required")
		return
	}

	err := handler.completeHandler.Handle(request.Context(), commands.CompleteBookingCommand{BookingID: bookingID})
	if err != nil {
		handleBookingError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, nil)
}

func (handler *BookingHandler) ListBookings(writer http.ResponseWriter, request *http.Request) {
	var dto dtos.CreateBookingDTO
	if err := json.NewDecoder(request.Body).Decode(&dto); err != nil {
		writeBadRequest(writer, "Invalid request body")
		return
	}

	if dto.UserID == "" && dto.GymID == "" {
		writeBadRequest(writer, "Either user_id or gym_id must be provided")
		return
	}

	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		writeBadRequest(writer, "Invalid start_time format")
		return
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		writeBadRequest(writer, "Invalid end_time format")
		return
	}

	result, err := handler.listHandler.Handle(request.Context(), queries.ListBookingsQuery{
		UserID:    dto.UserID,
		GymID:     dto.GymID,
		StartTime: startTime,
		EndTime:   endTime,
	})

	if err != nil {
		writeInternalError(writer)
		return
	}

	writeJSON(writer, http.StatusOK, result.Bookings)
}

func (handler *BookingHandler) CancelBooking(writer http.ResponseWriter, request *http.Request) {
	bookingID := request.PathValue("id")
	if bookingID == "" {
		writeBadRequest(writer, "Booking ID is required")
		return
	}

	err := handler.cancelHandler.Handle(request.Context(), commands.CancelBookingCommand{BookingID: bookingID})
	if err != nil {
		handleBookingError(writer, err)
		return
	}

	writeJSON(writer, http.StatusOK, nil)
}
