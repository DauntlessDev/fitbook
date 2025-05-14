package test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/yourusername/fitbook/booking-service/internal/application/commands"
	"github.com/yourusername/fitbook/booking-service/internal/application/dtos"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
	"github.com/yourusername/fitbook/booking-service/internal/domain/booking/test/mocks"
)

func TestCreateBookingHandler(t *testing.T) {
	repo := mocks.NewMockRepository()
	publisher := mocks.NewMockEventPublisher()
	handler := commands.NewCreateBookingHandler(repo, publisher)

	now := time.Now().Truncate(time.Second) // Truncate to seconds for consistent comparison
	validDTO := &dtos.CreateBookingDTO{
		UserID:    "user1",
		GymID:     "gym1",
		StartTime: now.Add(time.Hour).Format(time.RFC3339),
		EndTime:   now.Add(2 * time.Hour).Format(time.RFC3339),
	}

	tests := []struct {
		name    string
		dto     *dtos.CreateBookingDTO
		wantErr bool
	}{
		{
			name:    "valid booking",
			dto:     validDTO,
			wantErr: false,
		},
		{
			name: "invalid time range",
			dto: &dtos.CreateBookingDTO{
				UserID:    "user1",
				GymID:     "gym1",
				StartTime: now.Add(2 * time.Hour).Format(time.RFC3339),
				EndTime:   now.Add(time.Hour).Format(time.RFC3339),
			},
			wantErr: true,
		},
		{
			name: "past booking",
			dto: &dtos.CreateBookingDTO{
				UserID:    "user1",
				GymID:     "gym1",
				StartTime: now.Add(-time.Hour).Format(time.RFC3339),
				EndTime:   now.Add(time.Hour).Format(time.RFC3339),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := commands.CreateBookingCommand{
				DTO: tt.dto,
			}

			result, err := handler.Handle(context.Background(), cmd)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Booking)
				assert.Equal(t, tt.dto.UserID, result.Booking.UserID)
				assert.Equal(t, tt.dto.GymID, result.Booking.GymID)
				assert.Equal(t, tt.dto.StartTime, result.Booking.StartTime)
				assert.Equal(t, tt.dto.EndTime, result.Booking.EndTime)
				assert.Equal(t, "PENDING", result.Booking.Status) // Changed to uppercase to match domain model

				// Verify event was published
				events := publisher.GetEvents()
				assert.Len(t, events, 1)
				event := events[0]
				assert.Equal(t, "booking.created", event.EventName())
			}
		})
	}

	t.Run("overlapping booking", func(t *testing.T) {
		// Clear any previous bookings
		repo.Clear()
		publisher.Clear()

		// Create first booking
		cmd1 := commands.CreateBookingCommand{
			DTO: validDTO,
		}
		_, err := handler.Handle(context.Background(), cmd1)
		assert.NoError(t, err)

		// Try to create overlapping booking
		cmd2 := commands.CreateBookingCommand{
			DTO: &dtos.CreateBookingDTO{
				UserID:    "user2",
				GymID:     "gym1",                                                     // Same gym
				StartTime: now.Add(1*time.Hour + 15*time.Minute).Format(time.RFC3339), // Overlaps with first booking
				EndTime:   now.Add(2*time.Hour + 15*time.Minute).Format(time.RFC3339), // Overlaps with first booking
			},
		}
		result, err := handler.Handle(context.Background(), cmd2)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, booking.ErrOverlappingBooking)
	})
}
