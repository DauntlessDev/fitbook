package commands

import (
	"context"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type CreateBookingCommand struct {
	UserID    string    `json:"user_id" validate:"required"`
	GymID     string    `json:"gym_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

type CreateBookingHandler struct {
	bookingService *booking.Service
}

func NewCreateBookingHandler(bookingService *booking.Service) *CreateBookingHandler {
	return &CreateBookingHandler{
		bookingService: bookingService,
	}
}

func (h *CreateBookingHandler) Handle(ctx context.Context, cmd CreateBookingCommand) error {
	if err := validator.ValidateUserID(cmd.UserID); err != nil {
		return err
	}
	if err := validator.ValidateGymID(cmd.GymID); err != nil {
		return err
	}
	if err := validator.ValidateTimeRange(cmd.StartTime, cmd.EndTime); err != nil {
		return err
	}

	_, err := h.bookingService.CreateBooking(
		ctx,
		cmd.UserID,
		cmd.GymID,
		cmd.StartTime,
		cmd.EndTime,
	)
	return err
}
