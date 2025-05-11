package booking

import (
	"context"
	"time"
)

type Repository interface {
	Create(ctx context.Context, booking *Booking) error
	GetByID(ctx context.Context, id string) (*Booking, error)
	Update(ctx context.Context, booking *Booking) error
	DeleteByID(ctx context.Context, id string) error
	ListByUserID(ctx context.Context, userID string, startTime, endTime time.Time) ([]*Booking, error)
	ListByGymID(ctx context.Context, gymID string, startTime, endTime time.Time) ([]*Booking, error)
}
