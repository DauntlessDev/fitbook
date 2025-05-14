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
	repo      booking.Repository
	publisher booking.EventPublisher
}

func NewConfirmBookingHandler(repo booking.Repository, publisher booking.EventPublisher) *ConfirmBookingHandler {
	return &ConfirmBookingHandler{
		repo:      repo,
		publisher: publisher,
	}
}

func (h *ConfirmBookingHandler) Handle(ctx context.Context, cmd ConfirmBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	bookingRecord, err := h.repo.GetByID(ctx, cmd.BookingID)
	if err != nil {
		return err
	}

	if err := bookingRecord.Confirm(); err != nil {
		return err
	}

	if err := h.repo.Update(ctx, bookingRecord); err != nil {
		return err
	}

	event := booking.NewBookingEvent(bookingRecord, "confirmed")
	return h.publisher.Publish(event)
}
