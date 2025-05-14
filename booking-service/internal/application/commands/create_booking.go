package commands

import (
	"context"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type CreateBookingCommand struct {
	DTO *dtos.CreateBookingDTO
}

type CreateBookingResult struct {
	Booking *dtos.BookingDTO
}

type CreateBookingHandler struct {
	repo      booking.Repository
	publisher booking.EventPublisher
}

func NewCreateBookingHandler(repo booking.Repository, publisher booking.EventPublisher) *CreateBookingHandler {
	return &CreateBookingHandler{
		repo:      repo,
		publisher: publisher,
	}
}

func (h *CreateBookingHandler) Handle(ctx context.Context, cmd CreateBookingCommand) (*CreateBookingResult, error) {
	if err := validator.ValidateCreateBookingDTO(cmd.DTO); err != nil {
		return nil, err
	}

	userID, gymID, startTime, endTime, err := cmd.DTO.ToDomain()
	if err != nil {
		return nil, err
	}

	bookingRecord, err := booking.NewBooking(userID, gymID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	existingBookings, err := h.repo.ListByGymID(ctx, gymID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	for _, existing := range existingBookings {
		if bookingRecord.OverlapsWith(existing) {
			return nil, booking.ErrOverlappingBooking
		}
	}

	if err := h.repo.Create(ctx, bookingRecord); err != nil {
		return nil, err
	}

	event := booking.NewBookingEvent(bookingRecord, "created")
	if err := h.publisher.Publish(event); err != nil {
		// TODO: Consider implementing event publishing retry mechanism
		return nil, err
	}

	return &CreateBookingResult{
		Booking: dtos.FromDomain(bookingRecord),
	}, nil
}
