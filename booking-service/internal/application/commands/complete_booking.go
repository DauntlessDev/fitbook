package commands

import (
	"context"

	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type CompleteBookingCommand struct {
	BookingID string `json:"booking_id" validate:"required"`
}

type CompleteBookingHandler struct {
	bookingService *booking.Service
}

func NewCompleteBookingHandler(bookingService *booking.Service) *CompleteBookingHandler {
	return &CompleteBookingHandler{
		bookingService: bookingService,
	}
}

func (h *CompleteBookingHandler) Handle(ctx context.Context, cmd CompleteBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	return h.bookingService.CompleteBooking(ctx, cmd.BookingID)
}
