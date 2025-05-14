package validator

import (
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

func ValidateTimeRange(startTime, endTime time.Time) error {
	if startTime.IsZero() || endTime.IsZero() {
		return booking.ErrInvalidInput
	}
	if !startTime.Before(endTime) {
		return booking.ErrInvalidTimeRange
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

func ValidateCreateBookingDTO(dto *dtos.CreateBookingDTO) error {
	if err := ValidateUserID(dto.UserID); err != nil {
		return err
	}
	if err := ValidateGymID(dto.GymID); err != nil {
		return err
	}

	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		return booking.ErrInvalidInput
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		return booking.ErrInvalidInput
	}

	return ValidateTimeRange(startTime, endTime)
}

func ValidateListBookingsQuery(userID, gymID string, startTime, endTime time.Time) error {
	if userID != "" {
		if err := ValidateUserID(userID); err != nil {
			return err
		}
	}
	if gymID != "" {
		if err := ValidateGymID(gymID); err != nil {
			return err
		}
	}
	return ValidateTimeRange(startTime, endTime)
}
