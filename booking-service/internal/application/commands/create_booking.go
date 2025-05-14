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
	domainService *booking.Service
	publisher     booking.EventPublisher
}

func NewCreateBookingHandler(domainService *booking.Service, publisher booking.EventPublisher) *CreateBookingHandler {
	return &CreateBookingHandler{
		domainService: domainService,
		publisher:     publisher,
	}
}

func (h *CreateBookingHandler) Handle(ctx context.Context, cmd CreateBookingCommand) (*CreateBookingResult, error) {
	if err := validator.ValidateUserID(cmd.DTO.UserID); err != nil {
		return nil, err
	}
	if err := validator.ValidateGymID(cmd.DTO.GymID); err != nil {
		return nil, err
	}

	// Convert DTO times to domain times
	userID, gymID, startTime, endTime, err := cmd.DTO.ToDomain()
	if err != nil {
		return nil, err
	}

	if err := validator.ValidateTimeRange(startTime, endTime); err != nil {
		return nil, err
	}

	bookingRecord, err := h.domainService.CreateBooking(
		ctx,
		userID,
		gymID,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}

	// Publish event
	event := booking.NewBookingEvent(bookingRecord, "created")
	if err := h.publisher.Publish(event); err != nil {
		// TODO: Consider implementing event publishing retry mechanism
		return nil, err
	}

	return &CreateBookingResult{
		Booking: dtos.FromDomain(bookingRecord),
	}, nil
}
