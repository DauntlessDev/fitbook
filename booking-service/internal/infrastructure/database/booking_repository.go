package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/yourusername/fitbook/booking-service/internal/domain/booking"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{
		db: db,
	}
}

func (r *BookingRepository) Create(ctx context.Context, b *booking.Booking) error {
	query := `
		INSERT INTO bookings (id, user_id, gym_id, start_time, end_time, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		b.ID,
		b.UserID,
		b.GymID,
		b.StartTime,
		b.EndTime,
		b.Status,
		now,
		now,
	)
	return err
}

func (r *BookingRepository) GetByID(ctx context.Context, id string) (*booking.Booking, error) {
	query := `
		SELECT id, user_id, gym_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE id = $1
	`
	var b booking.Booking
	var createdAt, updatedAt time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID,
		&b.UserID,
		&b.GymID,
		&b.StartTime,
		&b.EndTime,
		&b.Status,
		&createdAt,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, booking.ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BookingRepository) Update(ctx context.Context, b *booking.Booking) error {
	query := `
		UPDATE bookings
		SET user_id = $1, gym_id = $2, start_time = $3, end_time = $4, status = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(ctx, query,
		b.UserID,
		b.GymID,
		b.StartTime,
		b.EndTime,
		b.Status,
		time.Now(),
		b.ID,
	)
	return err
}

func (r *BookingRepository) ListByUserID(ctx context.Context, userID string, startTime, endTime time.Time) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, gym_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE user_id = $1 AND start_time >= $2 AND end_time <= $3
		ORDER BY start_time ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*booking.Booking
	for rows.Next() {
		var b booking.Booking
		var createdAt, updatedAt time.Time
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.GymID,
			&b.StartTime,
			&b.EndTime,
			&b.Status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, rows.Err()
}

func (r *BookingRepository) ListByGymID(ctx context.Context, gymID string, startTime, endTime time.Time) ([]*booking.Booking, error) {
	query := `
		SELECT id, user_id, gym_id, start_time, end_time, status, created_at, updated_at
		FROM bookings
		WHERE gym_id = $1 AND start_time >= $2 AND end_time <= $3
		ORDER BY start_time ASC
	`
	rows, err := r.db.QueryContext(ctx, query, gymID, startTime, endTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []*booking.Booking
	for rows.Next() {
		var b booking.Booking
		var createdAt, updatedAt time.Time
		if err := rows.Scan(
			&b.ID,
			&b.UserID,
			&b.GymID,
			&b.StartTime,
			&b.EndTime,
			&b.Status,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, &b)
	}
	return bookings, rows.Err()
}

func (r *BookingRepository) DeleteByID(ctx context.Context, id string) error {
	query := `
		DELETE FROM bookings
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
