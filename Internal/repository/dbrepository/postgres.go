package dbrepo

import (
	"context"
	"time"

	"github.com/bharatsabne/bookings/Internal/models"
)

func (m *postgresDBRepo) AllUsers() bool {
	return true
}

// InsertReservation Inserts Reservation into database
func (m *postgresDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var newId int
	stament := `INSERT INTO Reservations
				(first_name, last_name, email_id, phone_number, start_date, end_date, 
				room_id, created_at, updated_at)
				VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := m.DB.QueryRowContext(ctx, stament,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomId,
		time.Now(),
		time.Now(),
	).Scan(&newId)
	if err != nil {
		return 0, err
	}
	return newId, nil
}

// InsertRoomRestrictions Inserts a room restriction into database
func (m *postgresDBRepo) InsertRoomRestrictions(r models.RoomRestriction) error {
	ctx, Cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer Cancel()
	stmt := `INSERT INTO public.room_restrictions
	(start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
	VALUES($1,$2,$3,$4,$5,$6,$7);`
	_, err := m.DB.ExecContext(ctx, stmt,
		r.StartDate,
		r.EndDate,
		r.RoomId,
		r.ReservationId,
		r.RestrictionId,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}
