package booking

import "time"

type Event interface {
	EventName() string
	OccurredAt() time.Time
}

type BaseBookingEvent struct {
	BookingID      string
	UserID         string
	GymID          string
	StartTime      time.Time
	EndTime        time.Time
	Status         BookingStatus
	OccurredAtTime time.Time
}

func (e BaseBookingEvent) OccurredAt() time.Time {
	return e.OccurredAtTime
}

type BookingCreatedEvent struct {
	BaseBookingEvent
}

func (e BookingCreatedEvent) EventName() string {
	return "booking.created"
}

type BookingCancelledEvent struct {
	BaseBookingEvent
}

func (e BookingCancelledEvent) EventName() string {
	return "booking.cancelled"
}

type BookingConfirmedEvent struct {
	BaseBookingEvent
}

func (e BookingConfirmedEvent) EventName() string {
	return "booking.confirmed"
}

type BookingCompletedEvent struct {
	BaseBookingEvent
}

func (e BookingCompletedEvent) EventName() string {
	return "booking.completed"
}
