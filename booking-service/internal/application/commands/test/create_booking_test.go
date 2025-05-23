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
	validBookingRequest := &dtos.CreateBookingDTO{
		UserID:    "user1",
		GymID:     "gym1",
		StartTime: now.Add(time.Hour).Format(time.RFC3339),
		EndTime:   now.Add(2 * time.Hour).Format(time.RFC3339),
	}

	tests := []struct {
		name              string
		request           *dtos.CreateBookingDTO
		shouldReturnError bool
	}{
		{
			name:              "valid booking",
			request:           validBookingRequest,
			shouldReturnError: false,
		},
		{
			name: "invalid time range",
			request: &dtos.CreateBookingDTO{
				UserID:    "user1",
				GymID:     "gym1",
				StartTime: now.Add(2 * time.Hour).Format(time.RFC3339),
				EndTime:   now.Add(time.Hour).Format(time.RFC3339),
			},
			shouldReturnError: true,
		},
		{
			name: "past booking",
			request: &dtos.CreateBookingDTO{
				UserID:    "user1",
				GymID:     "gym1",
				StartTime: now.Add(-time.Hour).Format(time.RFC3339),
				EndTime:   now.Add(time.Hour).Format(time.RFC3339),
			},
			shouldReturnError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := commands.CreateBookingCommand{
				DTO: test.request,
			}

			result, err := handler.Handle(context.Background(), command)
			if test.shouldReturnError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotNil(t, result.Booking)
				assert.Equal(t, test.request.UserID, result.Booking.UserID)
				assert.Equal(t, test.request.GymID, result.Booking.GymID)
				assert.Equal(t, test.request.StartTime, result.Booking.StartTime)
				assert.Equal(t, test.request.EndTime, result.Booking.EndTime)
				assert.Equal(t, "PENDING", result.Booking.Status)

				events := publisher.GetEvents()
				assert.Len(t, events, 1)
				event := events[0]
				assert.Equal(t, "booking.created", event.EventName())
			}
		})
	}

	t.Run("overlapping booking", func(t *testing.T) {
		repo.Clear()
		publisher.Clear()

		firstCommand := commands.CreateBookingCommand{
			DTO: validBookingRequest,
		}
		_, err := handler.Handle(context.Background(), firstCommand)
		assert.NoError(t, err)

		overlappingCommand := commands.CreateBookingCommand{
			DTO: &dtos.CreateBookingDTO{
				UserID:    "user2",
				GymID:     "gym1",
				StartTime: now.Add(1*time.Hour + 15*time.Minute).Format(time.RFC3339),
				EndTime:   now.Add(2*time.Hour + 15*time.Minute).Format(time.RFC3339),
			},
		}
		result, err := handler.Handle(context.Background(), overlappingCommand)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, booking.ErrOverlappingBooking)
	})
}
