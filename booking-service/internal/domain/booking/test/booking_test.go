package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

func TestNewBooking(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		gymID     string
		startTime time.Time
		endTime   time.Time
		wantErr   bool
	}{
		{
			name:      "valid booking",
			userID:    "user1",
			gymID:     "gym1",
			startTime: time.Now().Add(time.Hour),
			endTime:   time.Now().Add(2 * time.Hour),
			wantErr:   false,
		},
		{
			name:      "past booking",
			userID:    "user1",
			gymID:     "gym1",
			startTime: time.Now().Add(-time.Hour),
			endTime:   time.Now().Add(time.Hour),
			wantErr:   true,
		},
		{
			name:      "invalid time range",
			userID:    "user1",
			gymID:     "gym1",
			startTime: time.Now().Add(2 * time.Hour),
			endTime:   time.Now().Add(time.Hour),
			wantErr:   true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			newBooking, err := booking.NewBooking(testCase.userID, testCase.gymID, testCase.startTime, testCase.endTime)
			if testCase.wantErr {
				assert.Error(t, err)
				assert.Nil(t, newBooking)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, newBooking)
				assert.Equal(t, testCase.userID, newBooking.UserID)
				assert.Equal(t, testCase.gymID, newBooking.GymID)
				assert.Equal(t, testCase.startTime, newBooking.StartTime)
				assert.Equal(t, testCase.endTime, newBooking.EndTime)
				assert.Equal(t, booking.StatusPending, newBooking.Status)
			}
		})
	}
}

func TestBookingStatusTransitions(t *testing.T) {
	now := time.Now()
	testBooking, err := booking.NewBooking("user1", "gym1", now.Add(time.Hour), now.Add(2*time.Hour))
	assert.NoError(t, err)

	t.Run("valid status transitions", func(t *testing.T) {
		// Pending -> Confirmed
		err := testBooking.Confirm()
		assert.NoError(t, err)
		assert.Equal(t, booking.StatusConfirmed, testBooking.Status)

		// Confirmed -> Completed
		err = testBooking.Complete()
		assert.NoError(t, err)
		assert.Equal(t, booking.StatusCompleted, testBooking.Status)
	})

	t.Run("invalid status transitions", func(t *testing.T) {
		// Completed -> Confirmed (invalid)
		err := testBooking.Confirm()
		assert.Error(t, err)
		assert.Equal(t, booking.StatusCompleted, testBooking.Status)
	})

	t.Run("cancellation", func(t *testing.T) {
		cancellationTestBooking, err := booking.NewBooking("user1", "gym1", now.Add(time.Hour), now.Add(2*time.Hour))
		assert.NoError(t, err)

		// Cancel from Pending
		err = cancellationTestBooking.Cancel()
		assert.NoError(t, err)
		assert.Equal(t, booking.StatusCancelled, cancellationTestBooking.Status)

		// Try to cancel again
		err = cancellationTestBooking.Cancel()
		assert.Error(t, err)
		assert.Equal(t, booking.StatusCancelled, cancellationTestBooking.Status)
	})
}

func TestBookingOverlap(t *testing.T) {
	now := time.Now()
	baseBooking, err := booking.NewBooking("user1", "gym1", now.Add(time.Hour), now.Add(2*time.Hour))
	assert.NoError(t, err)

	tests := []struct {
		name        string
		startTime   time.Time
		endTime     time.Time
		wantOverlap bool
	}{
		{
			name:        "no overlap - before",
			startTime:   now.Add(30 * time.Minute),
			endTime:     now.Add(45 * time.Minute),
			wantOverlap: false,
		},
		{
			name:        "no overlap - after",
			startTime:   now.Add(2*time.Hour + time.Minute),
			endTime:     now.Add(2*time.Hour + 30*time.Minute),
			wantOverlap: false,
		},
		{
			name:        "overlap - starts during",
			startTime:   now.Add(1*time.Hour + 30*time.Minute),
			endTime:     now.Add(2*time.Hour + 30*time.Minute),
			wantOverlap: true,
		},
		{
			name:        "overlap - ends during",
			startTime:   now.Add(30 * time.Minute),
			endTime:     now.Add(1*time.Hour + 30*time.Minute),
			wantOverlap: true,
		},
		{
			name:        "overlap - completely contains",
			startTime:   now.Add(45 * time.Minute),
			endTime:     now.Add(2*time.Hour + 15*time.Minute),
			wantOverlap: true,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			comparisonBooking, err := booking.NewBooking("user1", "gym1", testCase.startTime, testCase.endTime)
			assert.NoError(t, err)
			assert.Equal(t, testCase.wantOverlap, baseBooking.OverlapsWith(comparisonBooking))
		})
	}
}
