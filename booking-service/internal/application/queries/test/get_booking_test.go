package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/fitbook/booking-service/internal/application/queries"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking/test/mocks"
)

func TestGetBookingHandler(t *testing.T) {
	repo := mocks.NewMockRepository()
	handler := queries.NewGetBookingHandler(repo)

	// Create a test booking
	now := time.Now().Truncate(time.Second) // Truncate to seconds
	testBooking, err := booking.NewBooking("user1", "gym1", now.Add(time.Hour), now.Add(2*time.Hour))
	assert.NoError(t, err)

	// Set a specific ID for testing
	testBooking.ID = "test-booking-1"
	repo.AddBooking(testBooking)

	tests := []struct {
		name      string
		bookingID string
		wantErr   bool
		errType   error
	}{
		{
			name:      "existing booking",
			bookingID: "test-booking-1",
			wantErr:   false,
		},
		{
			name:      "non-existent booking",
			bookingID: "non-existent",
			wantErr:   true,
			errType:   booking.ErrBookingNotFound,
		},
		{
			name:      "empty booking ID",
			bookingID: "",
			wantErr:   true,
			errType:   booking.ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			query := queries.GetBookingQuery{
				BookingID: tt.bookingID,
			}

			result, err := handler.Handle(context.Background(), query)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
				if tt.errType != nil {
					assert.ErrorIs(t, err, tt.errType)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Booking)
				assert.Equal(t, testBooking.ID, result.Booking.ID)
				assert.Equal(t, testBooking.UserID, result.Booking.UserID)
				assert.Equal(t, testBooking.GymID, result.Booking.GymID)

				// Parse the DTO time strings for comparison
				startTime, err := time.Parse(time.RFC3339, result.Booking.StartTime)
				assert.NoError(t, err)
				endTime, err := time.Parse(time.RFC3339, result.Booking.EndTime)
				assert.NoError(t, err)

				// Truncate times to seconds for comparison
				assert.Equal(t, testBooking.StartTime.Truncate(time.Second), startTime.Truncate(time.Second))
				assert.Equal(t, testBooking.EndTime.Truncate(time.Second), endTime.Truncate(time.Second))
				assert.Equal(t, testBooking.Status.String(), result.Booking.Status)
			}
		})
	}
}
