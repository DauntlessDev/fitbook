// All time fields use ISO8601 format (RFC3339) for consistency.
package dtos

import (
	"fmt"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type BookingDTO struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	GymID     string `json:"gym_id"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Status    string `json:"status"`
	Duration  int    `json:"duration"` // in minutes
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type CreateBookingDTO struct {
	UserID    string `json:"user_id" validate:"required"`
	GymID     string `json:"gym_id" validate:"required"`
	StartTime string `json:"start_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
	EndTime   string `json:"end_time" validate:"required,datetime=2006-01-02T15:04:05Z07:00"`
}

type UpdateBookingDTO struct {
	Status string `json:"status" validate:"required,oneof=PENDING CONFIRMED CANCELLED COMPLETED"`
}

type ErrorDTO struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

type ValidationErrorDTO struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

func (dto *CreateBookingDTO) ToDomain() (string, string, time.Time, time.Time, error) {
	startTime, err := time.Parse(time.RFC3339, dto.StartTime)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid start_time format: %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, dto.EndTime)
	if err != nil {
		return "", "", time.Time{}, time.Time{}, fmt.Errorf("invalid end_time format: %w", err)
	}

	return dto.UserID, dto.GymID, startTime, endTime, nil
}

func FromDomain(booking *booking.Booking) *BookingDTO {
	duration := int(booking.EndTime.Sub(booking.StartTime).Minutes())

	return &BookingDTO{
		ID:        booking.ID,
		UserID:    booking.UserID,
		GymID:     booking.GymID,
		StartTime: booking.StartTime.Format(time.RFC3339),
		EndTime:   booking.EndTime.Format(time.RFC3339),
		Status:    booking.Status.String(),
		Duration:  duration,
		CreatedAt: booking.CreatedAt.Format(time.RFC3339),
		UpdatedAt: booking.UpdatedAt.Format(time.RFC3339),
	}
}

func NewErrorDTO(code string, message string, details ...string) *ErrorDTO {
	return &ErrorDTO{
		Code:    code,
		Message: message,
		Details: details,
	}
}

func NewValidationErrorDTO(field string, message string, code string) *ValidationErrorDTO {
	return &ValidationErrorDTO{
		Field:   field,
		Message: message,
		Code:    code,
	}
}
