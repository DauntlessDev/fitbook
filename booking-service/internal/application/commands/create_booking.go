package commands

import (
	"context"
	"fmt"

	"github.com/google/uuid"
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

	existingBookings, err := h.repo.ListByGymID(ctx, gymID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	newBooking, err := booking.NewBooking(userID, gymID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	newBooking.ID = uuid.New().String()

	for _, existing := range existingBookings {
		if newBooking.OverlapsWith(existing) {
			return nil, booking.ErrOverlappingBooking
		}
	}

	fmt.Printf("new booking: %+v\n", newBooking)

	if err := h.repo.Create(ctx, newBooking); err != nil {
		return nil, err
	}

	event := booking.NewBookingEvent(newBooking, "created")
	if err := h.publisher.Publish(event); err != nil {
		// Log the error but don't fail the request
		// TODO: Add proper logging
	}

	return &CreateBookingResult{
		Booking: dtos.FromDomain(newBooking),
	}, nil
}
