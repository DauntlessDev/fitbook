package validator

import (
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

func ValidateTimeRange(startTime, endTime time.Time) error {
	if startTime.IsZero() || endTime.IsZero() {
		return booking.ErrInvalidInput
	}
	return nil
}

func ValidateRequiredString(value, fieldName string) error {
	if value == "" {
		return booking.ErrInvalidInput
	}
	return nil
}

func ValidateBookingID(bookingID string) error {
	return ValidateRequiredString(bookingID, "booking_id")
}

func ValidateUserID(userID string) error {
	return ValidateRequiredString(userID, "user_id")
}

func ValidateGymID(gymID string) error {
	return ValidateRequiredString(gymID, "gym_id")
}
