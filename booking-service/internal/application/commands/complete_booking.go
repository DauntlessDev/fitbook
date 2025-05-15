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
	repo      booking.Repository
	publisher booking.EventPublisher
}

func NewCompleteBookingHandler(repo booking.Repository, publisher booking.EventPublisher) *CompleteBookingHandler {
	return &CompleteBookingHandler{
		repo:      repo,
		publisher: publisher,
	}
}

func (handler *CompleteBookingHandler) Handle(ctx context.Context, cmd CompleteBookingCommand) error {
	if err := validator.ValidateBookingID(cmd.BookingID); err != nil {
		return err
	}

	bookingRecord, err := handler.repo.GetByID(ctx, cmd.BookingID)
	if err != nil {
		return err
	}

	if err := bookingRecord.Complete(); err != nil {
		return err
	}

	if err := handler.repo.Update(ctx, bookingRecord); err != nil {
		return err
	}

	event := booking.NewBookingEvent(bookingRecord, "completed")
	return handler.publisher.Publish(event)
}
