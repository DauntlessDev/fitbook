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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := booking.NewBooking(tt.userID, tt.gymID, tt.startTime, tt.endTime)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, b)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, b)
				assert.Equal(t, tt.userID, b.UserID)
				assert.Equal(t, tt.gymID, b.GymID)
				assert.Equal(t, tt.startTime, b.StartTime)
				assert.Equal(t, tt.endTime, b.EndTime)
				assert.Equal(t, booking.StatusPending, b.Status)
			}
		})
	}
}

func TestBookingStatusTransitions(t *testing.T) {
	now := time.Now()
	b, err := booking.NewBooking("user1", "gym1", now.Add(time.Hour), now.Add(2*time.Hour))
	assert.NoError(t, err)

	t.Run("valid status transitions", func(t *testing.T) {
		// Pending -> Confirmed
		err := b.Confirm()
		assert.NoError(t, err)
		assert.Equal(t, booking.StatusConfirmed, b.Status)

		// Confirmed -> Completed
		err = b.Complete()
		assert.NoError(t, err)
		assert.Equal(t, booking.StatusCompleted, b.Status)
	})

	t.Run("invalid status transitions", func(t *testing.T) {
		// Completed -> Confirmed (invalid)
		err := b.Confirm()
		assert.Error(t, err)
		assert.Equal(t, booking.StatusCompleted, b.Status)
	})

	t.Run("cancellation", func(t *testing.T) {
		b, err := booking.NewBooking("user1", "gym1", now.Add(time.Hour), now.Add(2*time.Hour))
		assert.NoError(t, err)

		// Cancel from Pending
		err = b.Cancel()
		assert.NoError(t, err)
		assert.Equal(t, booking.StatusCancelled, b.Status)

		// Try to cancel again
		err = b.Cancel()
		assert.Error(t, err)
		assert.Equal(t, booking.StatusCancelled, b.Status)
	})
}

func TestBookingOverlap(t *testing.T) {
	now := time.Now()
	b1, err := booking.NewBooking("user1", "gym1", now.Add(time.Hour), now.Add(2*time.Hour))
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b2, err := booking.NewBooking("user1", "gym1", tt.startTime, tt.endTime)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantOverlap, b1.OverlapsWith(b2))
		})
	}
}
