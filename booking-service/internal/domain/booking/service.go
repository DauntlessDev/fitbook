package booking

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateBooking(ctx context.Context, userID, gymID string, startTime, endTime time.Time) (*Booking, error) {
	booking, err := NewBooking(userID, gymID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	existingBookings, err := s.repo.ListByGymID(ctx, gymID, startTime, endTime)
	if err != nil {
		return nil, err
	}

	for _, existing := range existingBookings {
		if booking.OverlapsWith(existing) {
			return nil, ErrOverlappingBooking
		}
	}

	if err := s.repo.Create(ctx, booking); err != nil {
		return nil, err
	}

	return booking, nil
}

func (s *Service) CancelBooking(ctx context.Context, bookingID string) error {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := booking.Cancel(); err != nil {
		return err
	}

	return s.repo.Update(ctx, booking)
}

func (s *Service) ConfirmBooking(ctx context.Context, bookingID string) error {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := booking.Confirm(); err != nil {
		return err
	}

	return s.repo.Update(ctx, booking)
}

func (s *Service) CompleteBooking(ctx context.Context, bookingID string) error {
	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := booking.Complete(); err != nil {
		return err
	}

	return s.repo.Update(ctx, booking)
}

func (s *Service) GetBooking(ctx context.Context, bookingID string) (*Booking, error) {
	return s.repo.GetByID(ctx, bookingID)
}

func (s *Service) ListUserBookings(ctx context.Context, userID string, startTime, endTime time.Time) ([]*Booking, error) {
	return s.repo.ListByUserID(ctx, userID, startTime, endTime)
}

func (s *Service) ListGymBookings(ctx context.Context, gymID string, startTime, endTime time.Time) ([]*Booking, error) {
	return s.repo.ListByGymID(ctx, gymID, startTime, endTime)
}
