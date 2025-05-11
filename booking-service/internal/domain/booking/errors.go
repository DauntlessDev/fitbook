package booking

import "errors"

var (
	ErrBookingAlreadyCancelled = errors.New("booking is already cancelled")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrBookingNotFound         = errors.New("booking not found")
	ErrOverlappingBooking      = errors.New("booking overlaps with existing booking")
	ErrInvalidTimeRange        = errors.New("invalid time range")
	ErrPastBooking             = errors.New("cannot book in the past")
	ErrInvalidInput            = errors.New("invalid input")
)
