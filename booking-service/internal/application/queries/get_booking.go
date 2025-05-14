package queries

import (
	"context"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type GetBookingQuery struct {
	BookingID string `json:"booking_id" validate:"required"`
}

type GetBookingResult struct {
	Booking *dtos.BookingDTO
}

type GetBookingHandler struct {
	repo booking.Repository
}

func NewGetBookingHandler(repo booking.Repository) *GetBookingHandler {
	return &GetBookingHandler{
		repo: repo,
	}
}

func (h *GetBookingHandler) Handle(ctx context.Context, query GetBookingQuery) (*GetBookingResult, error) {
	if err := validator.ValidateBookingID(query.BookingID); err != nil {
		return nil, err
	}

	bookingRecord, err := h.repo.GetByID(ctx, query.BookingID)
	if err != nil {
		return nil, err
	}

	return &GetBookingResult{
		Booking: dtos.FromDomain(bookingRecord),
	}, nil
}
