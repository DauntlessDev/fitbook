package mocks

import (
	"sync"

	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type MockEventPublisher struct {
	mu     sync.RWMutex
	events []booking.Event
}

func NewMockEventPublisher() *MockEventPublisher {
	return &MockEventPublisher{
		events: make([]booking.Event, 0),
	}
}

func (m *MockEventPublisher) Publish(event booking.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = append(m.events, event)
	return nil
}

// Helper methods for testing
func (m *MockEventPublisher) GetEvents() []booking.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.events
}

func (m *MockEventPublisher) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = make([]booking.Event, 0)
}

func (m *MockEventPublisher) GetLastEvent() booking.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.events) == 0 {
		return nil
	}
	return m.events[len(m.events)-1]
}
