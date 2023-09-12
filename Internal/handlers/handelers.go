// Home Page Handler
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	helpers "github.com/bharatsabne/bookings/Internal/Helpers"
	"github.com/bharatsabne/bookings/Internal/config"
	driver "github.com/bharatsabne/bookings/Internal/drivers"
	"github.com/bharatsabne/bookings/Internal/forms"
	"github.com/bharatsabne/bookings/Internal/models"
	"github.com/bharatsabne/bookings/Internal/render"
	"github.com/bharatsabne/bookings/Internal/repository"
	dbrepo "github.com/bharatsabne/bookings/Internal/repository/dbrepository"
	"github.com/go-chi/chi/v5"
)

// Repo will be used by handlers
var Repo *Repositoy

type Repositoy struct {
	App *config.AppConfig
	DB  repository.Databaserepo
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repositoy {
	return &Repositoy{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandler Set the Repository for Handlers
func NewHandler(r *Repositoy) {
	Repo = r
}

func (m *Repositoy) Home(w http.ResponseWriter, r *http.Request) {
	// ip, _, err := net.SplitHostPort(r.RemoteAddr)
	// if err != nil {
	// 	//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	// 	fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
	// }

	// //fmt.Println("IP Address: " + ip)
	// m.App.Session.Put(r.Context(), "remote_ip", ip)
	render.Template(w, r, "Home.page.html", &models.TempateData{})
}

// About Page handler
func (m *Repositoy) About(w http.ResponseWriter, r *http.Request) {
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hi, There"
	// lol := m.App.Session.GetString(r.Context(), "remote_ip")
	// stringMap["remote_ip"] = lol
	// render.Template(w, r, "About.page.html", &models.TempateData{
	// 	StringMap: stringMap,
	// })
	render.Template(w, r, "About.page.html", &models.TempateData{})
}

// Contact
func (m *Repositoy) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "Contact.page.html", &models.TempateData{})
}

// Reservation
func (m *Repositoy) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cant get reservation from session"))
	}
	room, err := m.DB.GetRoomById(res.RoomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, r, "Make-Reservation.page.html", &models.TempateData{
		Forms:     forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// Post Reservation hanldes the posting of reservation form
func (m *Repositoy) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cant get reservation from session"))
	}
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation.FirstName = r.Form.Get("FirstName")
	reservation.LastName = r.Form.Get("LastName")
	reservation.Phone = r.Form.Get("Phone")
	reservation.Email = r.Form.Get("Email")

	form := forms.New(r.PostForm)

	form.Required("FirstName", "LastName", "Email")
	form.MinimumLength("FirstName", 3)
	form.MinimumLength("LastName", 3)
	form.IsEmail("Email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.Template(w, r, "Make-Reservation.page.html", &models.TempateData{
			Forms: form,
			Data:  data,
		})
		return
	}
	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	log.Println("newReservationId", newReservationId)
	restrictions := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}
	err = m.DB.InsertRoomRestrictions(restrictions)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "resarvation", reservation)
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals
func (m *Repositoy) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "Generals.page.html", &models.TempateData{})
}

// Marjors
func (m *Repositoy) Marjors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "Majors.page.html", &models.TempateData{})
}

// SearchAvailability
func (m *Repositoy) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "Search-Availability.page.html", &models.TempateData{})
}
func (m *Repositoy) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomId, err := strconv.Atoi(chi.URLParam(r, "Id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
	}
	res.RoomId = roomId
	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// POST SearchAvailability
func (m *Repositoy) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	rooms, err := m.DB.SearchForAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}
	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.html", &models.TempateData{Data: data})
	//w.Write([]byte(fmt.Sprintf("Start Date is %s and End Date is %s", start, end)))
}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	RoomId    string `json:"room_id"`
}

// AvailabilityJSON handles requests for avalibility and send JSON Response
func (m *Repositoy) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
	}
	available, err := m.DB.SearchForAvailabilityByDatesAndRoomId(startDate, endDate, roomId)
	if err != nil {
		helpers.ServerError(w, err)
	}
	rest := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomId),
	}
	out, err := json.MarshalIndent(rest, "", "     ")
	if err != nil {
		// log.Println(err)
		helpers.ServerError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// ReservationSummary Displayes Reservation Summary page
func (m *Repositoy) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	resarvation, ok := m.App.Session.Get(r.Context(), "resarvation").(models.Reservation)
	if !ok {
		//log.Println("Cannot get session")
		m.App.ErrorLog.Println("can not get error from session")
		m.App.Session.Put(r.Context(), "error", "can't get reservation form session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	m.App.Session.Remove(r.Context(), "resarvation")

	data := make(map[string]interface{})
	data["resarvation"] = resarvation

	sd := resarvation.StartDate.Format("2006-01-02")
	ed := resarvation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	render.Template(w, r, "Reservation-Summery.page.html", &models.TempateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// BookRoom takes url param and builds session and takes user to reservation page
func (m *Repositoy) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
	}
	room, err := m.DB.GetRoomById(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var res models.Reservation
	res.RoomId = roomID
	res.Room.RoomName = room.RoomName
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}
