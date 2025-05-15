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
	repo      booking.Repository
	publisher booking.EventPublisher
}

func NewCancelBookingHandler(repo booking.Repository, publisher booking.EventPublisher) *CancelBookingHandler {
	return &CancelBookingHandler{
		repo:      repo,
		publisher: publisher,
	}
}

func (handler *CancelBookingHandler) Handle(ctx context.Context, cmd CancelBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	bookingRecord, err := handler.repo.GetByID(ctx, cmd.BookingID)
	if err != nil {
		return err
	}

	if err := bookingRecord.Cancel(); err != nil {
		return err
	}

	if err := handler.repo.Update(ctx, bookingRecord); err != nil {
		return err
	}

	event := booking.NewBookingEvent(bookingRecord, "cancelled")
	return handler.publisher.Publish(event)
}
