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
	bookingService *booking.Service
}

func NewCreateBookingHandler(bookingService *booking.Service) *CreateBookingHandler {
	return &CreateBookingHandler{
		bookingService: bookingService,
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

	booking, err := h.bookingService.CreateBooking(
		ctx,
		userID,
		gymID,
		startTime,
		endTime,
	)
	if err != nil {
		return nil, err
	}

	return &CreateBookingResult{
		Booking: dtos.FromDomain(booking),
	}, nil
}
