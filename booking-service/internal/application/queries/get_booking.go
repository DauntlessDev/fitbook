package queries

import (
	"context"

	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type GetBookingQuery struct {
	BookingID string `json:"booking_id" validate:"required"`
}

type GetBookingResult struct {
	Booking *booking.Booking
}

type GetBookingHandler struct {
	bookingService *booking.Service
}

func NewGetBookingHandler(bookingService *booking.Service) *GetBookingHandler {
	return &GetBookingHandler{
		bookingService: bookingService,
	}
}

func (h *GetBookingHandler) Handle(ctx context.Context, query GetBookingQuery) (*GetBookingResult, error) {
	if err := validator.ValidateBookingID(query.BookingID); err != nil {
		return nil, err
	}

	booking, err := h.bookingService.GetBooking(ctx, query.BookingID)
	if err != nil {
		return nil, err
	}

	return &GetBookingResult{
		Booking: booking,
	}, nil
}
