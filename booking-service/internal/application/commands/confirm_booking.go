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
	domainService *booking.Service
	publisher     booking.EventPublisher
}

func NewConfirmBookingHandler(domainService *booking.Service, publisher booking.EventPublisher) *ConfirmBookingHandler {
	return &ConfirmBookingHandler{
		domainService: domainService,
		publisher:     publisher,
	}
}

func (h *ConfirmBookingHandler) Handle(ctx context.Context, cmd ConfirmBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	if err := h.domainService.ConfirmBooking(ctx, cmd.BookingID); err != nil {
		return err
	}

	bookingRecord, err := h.domainService.GetBooking(ctx, cmd.BookingID)
	if err != nil {
		return err
	}

	event := booking.NewBookingEvent(bookingRecord, "confirmed")
	return h.publisher.Publish(event)
}
