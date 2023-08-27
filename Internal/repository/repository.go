package repository

import "github.com/bharatsabne/bookings/Internal/models"

type Databaserepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestrictions(r models.RoomRestriction) error
}
