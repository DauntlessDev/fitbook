package booking

type BookingStatus string

const (
	StatusPending   BookingStatus = "PENDING"
	StatusConfirmed BookingStatus = "CONFIRMED"
	StatusCancelled BookingStatus = "CANCELLED"
	StatusCompleted BookingStatus = "COMPLETED"
)

func (s BookingStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusConfirmed, StatusCancelled, StatusCompleted:
		return true
	default:
		return false
	}
}

func (s BookingStatus) String() string {
	return string(s)
}
