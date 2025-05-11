package commands

import (
	"context"

	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type CancelBookingCommand struct {
	BookingID string `json:"booking_id" validate:"required"`
}

type CancelBookingHandler struct {
	bookingService *booking.Service
}

func NewCancelBookingHandler(bookingService *booking.Service) *CancelBookingHandler {
	return &CancelBookingHandler{
		bookingService: bookingService,
	}
}

func (h *CancelBookingHandler) Handle(ctx context.Context, cmd CancelBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	return h.bookingService.CancelBooking(ctx, cmd.BookingID)
}
