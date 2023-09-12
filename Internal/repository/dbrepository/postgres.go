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

// SearchForAvailabilityByDates Check for availability by dates return true when Availble else false for room id
func (m *postgresDBRepo) SearchForAvailabilityByDatesAndRoomId(start, end time.Time, roomId int) (bool, error) {
	ctx, Cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer Cancel()
	stmt := `SELECT COUNT(id) 
			FROM room_restrictions 
			WHERE room_id= $1 and $2 < end_date and $3 > start_date`
	var numRows int
	row := m.DB.QueryRowContext(ctx, stmt, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}
	if numRows == 0 {
		return true, nil
	}
	return false, nil
}

// SearchForAvailabilityForAllRooms returns slice wo rooms if any from selected date range
func (m *postgresDBRepo) SearchForAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, Cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer Cancel()
	var rooms []models.Room
	stmt := `select r.id, r.room_name from rooms r 
			where r.id not in (select rr.room_id from room_restrictions rr 
			where $1 <= rr.end_date and $2 >= rr.start_date)`
	rows, err := m.DB.QueryContext(ctx, stmt, start, end)
	if err != nil {
		return rooms, err
	}
	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(
			&room.Id,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}
	if err := rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomById Gets room by Id
func (m *postgresDBRepo) GetRoomById(id int) (models.Room, error) {
	ctx, Cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer Cancel()
	var room models.Room
	query := `select id, room_name, created_at, updated_at from rooms where id=$1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.Id,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, nil
	}
	return room, nil
}
