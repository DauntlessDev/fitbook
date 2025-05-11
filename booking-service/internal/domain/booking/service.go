package booking

import (
	"context"
	"time"
)

type Service struct {
	repo      Repository
	publisher EventPublisher
}

func NewService(repo Repository, publisher EventPublisher) *Service {
	return &Service{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *Service) CreateBooking(ctx context.Context, userID, gymID string, startTime, endTime time.Time) (*Booking, error) {
	if userID == "" || gymID == "" {
		return nil, ErrInvalidInput
	}
	if startTime.IsZero() || endTime.IsZero() {
		return nil, ErrInvalidInput
	}

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

	event := NewBookingEvent(booking, "created")
	if err := s.publisher.Publish(event); err != nil {
		// TODO: Consider implementing event publishing retry mechanism
		// For now, we'll just log the error
		return nil, err
	}

	return booking, nil
}

func (s *Service) CancelBooking(ctx context.Context, bookingID string) error {
	if bookingID == "" {
		return ErrInvalidInput
	}

	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := booking.Cancel(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, booking); err != nil {
		return err
	}

	event := NewBookingEvent(booking, "cancelled")
	return s.publisher.Publish(event)
}

func (s *Service) ConfirmBooking(ctx context.Context, bookingID string) error {
	if bookingID == "" {
		return ErrInvalidInput
	}

	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := booking.Confirm(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, booking); err != nil {
		return err
	}

	event := NewBookingEvent(booking, "confirmed")
	return s.publisher.Publish(event)
}

func (s *Service) CompleteBooking(ctx context.Context, bookingID string) error {
	if bookingID == "" {
		return ErrInvalidInput
	}

	booking, err := s.repo.GetByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if err := booking.Complete(); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, booking); err != nil {
		return err
	}

	event := NewBookingEvent(booking, "completed")
	return s.publisher.Publish(event)
}

func (s *Service) GetBooking(ctx context.Context, bookingID string) (*Booking, error) {
	if bookingID == "" {
		return nil, ErrInvalidInput
	}
	return s.repo.GetByID(ctx, bookingID)
}

func (s *Service) ListUserBookings(ctx context.Context, userID string, startTime, endTime time.Time) ([]*Booking, error) {
	if userID == "" || startTime.IsZero() || endTime.IsZero() {
		return nil, ErrInvalidInput
	}
	return s.repo.ListByUserID(ctx, userID, startTime, endTime)
}

func (s *Service) ListGymBookings(ctx context.Context, gymID string, startTime, endTime time.Time) ([]*Booking, error) {
	if gymID == "" || startTime.IsZero() || endTime.IsZero() {
		return nil, ErrInvalidInput
	}
	return s.repo.ListByGymID(ctx, gymID, startTime, endTime)
}
