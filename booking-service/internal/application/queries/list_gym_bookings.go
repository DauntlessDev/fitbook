package queries

import (
	"context"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type ListGymBookingsQuery struct {
	GymID     string    `json:"gym_id" validate:"required"`
	StartTime time.Time `json:"start_time" validate:"required"`
	EndTime   time.Time `json:"end_time" validate:"required"`
}

type ListGymBookingsResult struct {
	Bookings []*dtos.BookingDTO
}

type ListGymBookingsHandler struct {
	repo booking.Repository
}

func NewListGymBookingsHandler(repo booking.Repository) *ListGymBookingsHandler {
	return &ListGymBookingsHandler{
		repo: repo,
	}
}

func (h *ListGymBookingsHandler) Handle(ctx context.Context, query ListGymBookingsQuery) (*ListGymBookingsResult, error) {
	if err := validator.ValidateGymID(query.GymID); err != nil {
		return nil, err
	}
	if err := validator.ValidateTimeRange(query.StartTime, query.EndTime); err != nil {
		return nil, err
	}

	bookings, err := h.repo.ListByGymID(ctx, query.GymID, query.StartTime, query.EndTime)
	if err != nil {
		return nil, err
	}

	bookingDTOs := make([]*dtos.BookingDTO, len(bookings))
	for i, bookingRecord := range bookings {
		bookingDTOs[i] = dtos.FromDomain(bookingRecord)
	}

	return &ListGymBookingsResult{
		Bookings: bookingDTOs,
	}, nil
}
