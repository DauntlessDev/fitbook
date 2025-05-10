package booking

import "errors"

var (
	ErrBookingAlreadyCancelled = errors.New("booking is already cancelled")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrInvalidTimeRange        = errors.New("invalid time range")
)
