package models

import "time"


//User is the user modal
type User struct {
	Id          int
	FirstName   string
	LastName    string
	Email       string
	Password    string
	AccessLevel int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

//Room is the room modal
type Room struct {
	Id        int
	RoomName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

//Restriction is Restriction modal
type Restriction struct {
	Id              int
	RestrictionName string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

//Reservation is the Reservation modal
type Reservation struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
	StartDate time.Time
	EndDate   time.Time
	RoomId    int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room      Room
}

//RoomRestriction is RoomRestriction modal
type RoomRestriction struct {
	Id            int
	StartDate     time.Time
	EndDate       time.Time
	RoomId        int
	ReservationId int
	RestrictionId int
	Room          Room
	Reservation   Reservation
	Restriction   Restriction
}
