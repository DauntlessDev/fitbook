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

func (m *MockRepository) Create(ctx context.Context, booking *booking.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bookings[booking.ID] = booking
	return nil
}

func (m *MockRepository) Update(ctx context.Context, booking *booking.Booking) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bookings[booking.ID] = booking
	return nil
}

func (m *MockRepository) GetByID(ctx context.Context, id string) (*booking.Booking, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if booking, exists := m.bookings[id]; exists {
		return booking, nil
	}
	return nil, booking.ErrBookingNotFound
}

func (m *MockRepository) ListByUserID(ctx context.Context, userID string, startTime, endTime time.Time) ([]*booking.Booking, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*booking.Booking
	for _, b := range m.bookings {
		if b.UserID == userID && b.StartTime.After(startTime) && b.EndTime.Before(endTime) {
			result = append(result, b)
		}
	}
	return result, nil
}

func (m *MockRepository) ListByGymID(ctx context.Context, gymID string, startTime, endTime time.Time) ([]*booking.Booking, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*booking.Booking
	for _, b := range m.bookings {
		if b.GymID == gymID {
			if b.StartTime.Before(endTime) && b.EndTime.After(startTime) {
				result = append(result, b)
			}
		}
	}
	return result, nil
}

// Helper methods for testing
func (m *MockRepository) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bookings = make(map[string]*booking.Booking)
}

func (m *MockRepository) AddBooking(b *booking.Booking) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.bookings[b.ID] = b
}

func (m *MockRepository) DeleteByID(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.bookings, id)
	return nil
}
