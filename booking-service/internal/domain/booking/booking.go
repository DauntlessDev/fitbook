package booking

import (
	"time"
)

type Booking struct {
	ID        string
	UserID    string
	GymID     string
	StartTime time.Time
	EndTime   time.Time
	Status    BookingStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBooking(userID, gymID string, startTime, endTime time.Time) (*Booking, error) {
	now := time.Now()

	if startTime.Before(now) {
		return nil, ErrPastBooking
	}

	if !startTime.Before(endTime) {
		return nil, ErrInvalidTimeRange
	}

	booking := &Booking{
		UserID:    userID,
		GymID:     gymID,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return booking, nil
}

func (booking *Booking) Cancel() error {
	if booking.Status == StatusCancelled {
		return ErrBookingAlreadyCancelled
	}
	if booking.Status == StatusCompleted {
		return ErrInvalidStatusTransition
	}
	booking.Status = StatusCancelled
	booking.UpdatedAt = time.Now()
	return nil
}

func (booking *Booking) Confirm() error {
	if booking.Status != StatusPending {
		return ErrInvalidStatusTransition
	}
	booking.Status = StatusConfirmed
	booking.UpdatedAt = time.Now()
	return nil
}

func (booking *Booking) Complete() error {
	if booking.Status != StatusConfirmed {
		return ErrInvalidStatusTransition
	}
	booking.Status = StatusCompleted
	booking.UpdatedAt = time.Now()
	return nil
}

func (booking *Booking) OverlapsWith(other *Booking) bool {
	return booking.GymID == other.GymID &&
		booking.StartTime.Before(other.EndTime) &&
		booking.EndTime.After(other.StartTime)
}
