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
	domainService *booking.Service
	publisher     booking.EventPublisher
}

func NewCancelBookingHandler(domainService *booking.Service, publisher booking.EventPublisher) *CancelBookingHandler {
	return &CancelBookingHandler{
		domainService: domainService,
		publisher:     publisher,
	}
}

func (h *CancelBookingHandler) Handle(ctx context.Context, cmd CancelBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	if err := h.domainService.CancelBooking(ctx, cmd.BookingID); err != nil {
		return err
	}

	bookingRecord, err := h.domainService.GetBooking(ctx, cmd.BookingID)
	if err != nil {
		return err
	}

	event := booking.NewBookingEvent(bookingRecord, "cancelled")
	return h.publisher.Publish(event)
}
