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

func NewBooking(userID, gymID string, startTime, endTime time.Time) *Booking {
	now := time.Now()
	return &Booking{
		UserID:    userID,
		GymID:     gymID,
		StartTime: startTime,
		EndTime:   endTime,
		Status:    StatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (b *Booking) Cancel() error {
	if b.Status == StatusCancelled {
		return ErrBookingAlreadyCancelled
	}
	b.Status = StatusCancelled
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Booking) Confirm() error {
	if b.Status != StatusPending {
		return ErrInvalidStatusTransition
	}
	b.Status = StatusConfirmed
	b.UpdatedAt = time.Now()
	return nil
}

func (b *Booking) Complete() error {
	if b.Status != StatusConfirmed {
		return ErrInvalidStatusTransition
	}
	b.Status = StatusCompleted
	b.UpdatedAt = time.Now()
	return nil
}
