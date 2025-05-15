package mocks

import (
	"context"
	"sync"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type MockRepository struct {
	mu       sync.RWMutex
	bookings map[string]*booking.Booking
}

func NewMockRepository() *MockRepository {
	return &MockRepository{
		bookings: make(map[string]*booking.Booking),
	}
}

func (repo *MockRepository) Create(ctx context.Context, booking *booking.Booking) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.bookings[booking.ID] = booking

	return nil
}

func (repo *MockRepository) Update(ctx context.Context, booking *booking.Booking) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.bookings[booking.ID] = booking

	return nil
}

func (repo *MockRepository) GetByID(ctx context.Context, id string) (*booking.Booking, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	if booking, exists := repo.bookings[id]; exists {
		return booking, nil
	}

	return nil, booking.ErrBookingNotFound
}

func (repo *MockRepository) ListByUserID(ctx context.Context, userID string, startTime, endTime time.Time) ([]*booking.Booking, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var result []*booking.Booking
	for _, booking := range repo.bookings {
		if booking.UserID == userID && booking.StartTime.After(startTime) && booking.EndTime.Before(endTime) {
			result = append(result, booking)
		}
	}

	return result, nil
}

func (repo *MockRepository) ListByGymID(ctx context.Context, gymID string, startTime, endTime time.Time) ([]*booking.Booking, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	var result []*booking.Booking
	for _, booking := range repo.bookings {
		if booking.GymID == gymID {
			if booking.StartTime.Before(endTime) && booking.EndTime.After(startTime) {
				result = append(result, booking)
			}
		}
	}

	return result, nil
}

func (repo *MockRepository) Clear() {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.bookings = make(map[string]*booking.Booking)
}

func (repo *MockRepository) AddBooking(booking *booking.Booking) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	repo.bookings[booking.ID] = booking
}

func (repo *MockRepository) DeleteByID(ctx context.Context, id string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	delete(repo.bookings, id)

	return nil
}
