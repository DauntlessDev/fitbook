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
	domainService *booking.Service
	publisher     booking.EventPublisher
}

func NewCompleteBookingHandler(domainService *booking.Service, publisher booking.EventPublisher) *CompleteBookingHandler {
	return &CompleteBookingHandler{
		domainService: domainService,
		publisher:     publisher,
	}
}

func (h *CompleteBookingHandler) Handle(ctx context.Context, cmd CompleteBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	if err := h.domainService.CompleteBooking(ctx, cmd.BookingID); err != nil {
		return err
	}

	bookingRecord, err := h.domainService.GetBooking(ctx, cmd.BookingID)
	if err != nil {
		return err
	}

	event := booking.NewBookingEvent(bookingRecord, "completed")
	return h.publisher.Publish(event)
}
