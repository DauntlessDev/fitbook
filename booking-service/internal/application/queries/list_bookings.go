package queries

import (
	"context"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/application/validator"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type ListBookingsQuery struct {
	UserID    string    `json:"user_id"`
	GymID     string    `json:"gym_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

type ListBookingsResult struct {
	Bookings []*dtos.BookingDTO
}

type ListBookingsHandler struct {
	repo booking.Repository
}

func NewListBookingsHandler(repo booking.Repository) *ListBookingsHandler {
	return &ListBookingsHandler{
		repo: repo,
	}
}

func (h *ListBookingsHandler) Handle(ctx context.Context, query ListBookingsQuery) (*ListBookingsResult, error) {
	if err := validator.ValidateListBookingsQuery(query.UserID, query.GymID, query.StartTime, query.EndTime); err != nil {
		return nil, err
	}

	var bookings []*booking.Booking
	var err error

	if query.UserID != "" {
		bookings, err = h.repo.ListByUserID(ctx, query.UserID, query.StartTime, query.EndTime)
	} else if query.GymID != "" {
		bookings, err = h.repo.ListByGymID(ctx, query.GymID, query.StartTime, query.EndTime)
	} else {
		return nil, booking.ErrInvalidInput
	}

	if err != nil {
		return nil, err
	}

	result := &ListBookingsResult{
		Bookings: make([]*dtos.BookingDTO, len(bookings)),
	}

	for i, b := range bookings {
		result.Bookings[i] = dtos.FromDomain(b)
	}

	return result, nil
}
