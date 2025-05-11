package commands

import (
	"context"

	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type ConfirmBookingCommand struct {
	BookingID string `json:"booking_id" validate:"required"`
}

type ConfirmBookingHandler struct {
	bookingService *booking.Service
}

func NewConfirmBookingHandler(bookingService *booking.Service) *ConfirmBookingHandler {
	return &ConfirmBookingHandler{
		bookingService: bookingService,
	}
}

func (h *ConfirmBookingHandler) Handle(ctx context.Context, cmd ConfirmBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	return h.bookingService.ConfirmBooking(ctx, cmd.BookingID)
}
