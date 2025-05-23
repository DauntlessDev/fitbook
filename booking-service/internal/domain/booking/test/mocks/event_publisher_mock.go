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

func (publisher *MockEventPublisher) Publish(event booking.Event) error {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	publisher.events = append(publisher.events, event)

	return nil
}

func (publisher *MockEventPublisher) GetEvents() []booking.Event {
	publisher.mu.RLock()
	defer publisher.mu.RUnlock()

	return publisher.events
}

func (publisher *MockEventPublisher) Clear() {
	publisher.mu.Lock()
	defer publisher.mu.Unlock()

	publisher.events = make([]booking.Event, 0)
}

func (publisher *MockEventPublisher) GetLastEvent() booking.Event {
	publisher.mu.RLock()
	defer publisher.mu.RUnlock()

	if len(publisher.events) == 0 {
		return nil
	}

	return publisher.events[len(publisher.events)-1]
}
