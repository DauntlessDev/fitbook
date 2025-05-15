package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dtos.Response{Success: true, Data: data})
}

func writeError(w http.ResponseWriter, status int, code, message string, details ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dtos.Response{
		Success: false,
		Error:   dtos.NewErrorDTO(code, message, details...),
	})
}

func writeBadRequest(w http.ResponseWriter, message string, details ...string) {
	writeError(w, http.StatusBadRequest, "INVALID_REQUEST", message, details...)
}

func writeInternalError(w http.ResponseWriter) {
	writeError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error")
}

func handleBookingError(w http.ResponseWriter, err error) {
	switch err {
	case booking.ErrInvalidTimeRange:
		writeError(w, http.StatusBadRequest, "INVALID_TIME_RANGE", err.Error())
	case booking.ErrPastBooking:
		writeError(w, http.StatusBadRequest, "PAST_BOOKING", err.Error())
	case booking.ErrOverlappingBooking:
		writeError(w, http.StatusConflict, "OVERLAPPING_BOOKING", err.Error())
	case booking.ErrBookingNotFound:
		writeError(w, http.StatusNotFound, "BOOKING_NOT_FOUND", err.Error())
	case booking.ErrBookingAlreadyCancelled:
		writeError(w, http.StatusBadRequest, "BOOKING_ALREADY_CANCELLED", err.Error())
	case booking.ErrInvalidStatusTransition:
		writeError(w, http.StatusBadRequest, "INVALID_STATUS_TRANSITION", err.Error())
	case booking.ErrInvalidInput:
		writeError(w, http.StatusBadRequest, "INVALID_INPUT", err.Error())
	default:
		writeInternalError(w)
	}
}
