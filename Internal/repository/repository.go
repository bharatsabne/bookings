package repository

import (
	"time"

	"github.com/bharatsabne/bookings/Internal/models"
)

type Databaserepo interface {
	AllUsers() bool
	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestrictions(r models.RoomRestriction) error
	SearchForAvailabilityByDatesAndRoomId(start, end time.Time, roomId int) (bool, error)
	SearchForAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomById(id int) (models.Room, error)
}
