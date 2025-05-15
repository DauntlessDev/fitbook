package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

func writeJSON(writer http.ResponseWriter, status int, data interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(dtos.Response{Success: true, Data: data})
}

func writeError(writer http.ResponseWriter, status int, code, message string, details ...string) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(dtos.Response{
		Success: false,
		Error:   dtos.NewErrorDTO(code, message, details...),
	})
}

func writeBadRequest(writer http.ResponseWriter, message string, details ...string) {
	writeError(writer, http.StatusBadRequest, "INVALID_REQUEST", message, details...)
}

func writeInternalError(writer http.ResponseWriter) {
	writeError(writer, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
}

func handleBookingError(writer http.ResponseWriter, err error) {
	switch err {
	case booking.ErrInvalidTimeRange:
		writeError(writer, http.StatusBadRequest, "INVALID_TIME_RANGE", err.Error())
	case booking.ErrPastBooking:
		writeError(writer, http.StatusBadRequest, "PAST_BOOKING", err.Error())
	case booking.ErrOverlappingBooking:
		writeError(writer, http.StatusConflict, "OVERLAPPING_BOOKING", err.Error())
	case booking.ErrBookingNotFound:
		writeError(writer, http.StatusNotFound, "BOOKING_NOT_FOUND", err.Error())
	case booking.ErrBookingAlreadyCancelled:
		writeError(writer, http.StatusBadRequest, "BOOKING_ALREADY_CANCELLED", err.Error())
	case booking.ErrInvalidStatusTransition:
		writeError(writer, http.StatusBadRequest, "INVALID_STATUS_TRANSITION", err.Error())
	case booking.ErrInvalidInput:
		writeError(writer, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	default:
		writeInternalError(writer)
	}
}
