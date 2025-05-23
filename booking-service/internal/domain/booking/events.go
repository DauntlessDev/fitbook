package booking

import "time"

type Event interface {
	EventName() string
	OccurredAt() time.Time
}

type EventPublisher interface {
	Publish(event Event) error
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

func (event BaseBookingEvent) OccurredAt() time.Time {
	return event.OccurredAtTime
}

type BookingCreatedEvent struct {
	BaseBookingEvent
}

func (event BookingCreatedEvent) EventName() string {
	return "booking.created"
}

type BookingCancelledEvent struct {
	BaseBookingEvent
}

func (event BookingCancelledEvent) EventName() string {
	return "booking.cancelled"
}

type BookingConfirmedEvent struct {
	BaseBookingEvent
}

func (event BookingConfirmedEvent) EventName() string {
	return "booking.confirmed"
}

type BookingCompletedEvent struct {
	BaseBookingEvent
}

func (event BookingCompletedEvent) EventName() string {
	return "booking.completed"
}

func NewBookingEvent(booking *Booking, eventType string) Event {
	baseEvent := BaseBookingEvent{
		BookingID:      booking.ID,
		UserID:         booking.UserID,
		GymID:          booking.GymID,
		StartTime:      booking.StartTime,
		EndTime:        booking.EndTime,
		Status:         booking.Status,
		OccurredAtTime: time.Now(),
	}

	switch eventType {
	case "created":
		return BookingCreatedEvent{BaseBookingEvent: baseEvent}
	case "cancelled":
		return BookingCancelledEvent{BaseBookingEvent: baseEvent}
	case "confirmed":
		return BookingConfirmedEvent{BaseBookingEvent: baseEvent}
	case "completed":
		return BookingCompletedEvent{BaseBookingEvent: baseEvent}
	default:
		return nil
	}
}
