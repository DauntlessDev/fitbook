package queries

import (
	"context"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type ListUserBookingsQuery struct {
	UserID    string    `json:"user_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

type ListUserBookingsResult struct {
	Bookings []*dtos.BookingDTO
}

type ListUserBookingsHandler struct {
	bookingService *booking.Service
}

func NewListUserBookingsHandler(bookingService *booking.Service) *ListUserBookingsHandler {
	return &ListUserBookingsHandler{
		bookingService: bookingService,
	}
}

func (h *ListUserBookingsHandler) Handle(ctx context.Context, query ListUserBookingsQuery) (*ListUserBookingsResult, error) {
	if err := validator.ValidateUserID(query.UserID); err != nil {
		return nil, err
	}
	if err := validator.ValidateTimeRange(query.StartTime, query.EndTime); err != nil {
		return nil, err
	}

	bookings, err := h.bookingService.ListUserBookings(ctx, query.UserID, query.StartTime, query.EndTime)
	if err != nil {
		return nil, err
	}

	bookingDTOs := make([]*dtos.BookingDTO, len(bookings))
	for i, b := range bookings {
		bookingDTOs[i] = dtos.FromDomain(b)
	}

	return &ListUserBookingsResult{
		Bookings: bookingDTOs,
	}, nil
}
