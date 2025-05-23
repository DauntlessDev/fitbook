package booking

type BookingStatus string

const (
	StatusPending   BookingStatus = "PENDING"
	StatusConfirmed BookingStatus = "CONFIRMED"
	StatusCancelled BookingStatus = "CANCELLED"
	StatusCompleted BookingStatus = "COMPLETED"
)

func (status BookingStatus) IsValid() bool {
	switch status {
	case StatusPending, StatusConfirmed, StatusCancelled, StatusCompleted:
		return true
	default:
		return false
	}
}

func (status BookingStatus) String() string {
	return string(status)
}
