// Home Page Handler
package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/bharatsabne/bookings/Internal/config"
	"github.com/bharatsabne/bookings/Internal/forms"
	"github.com/bharatsabne/bookings/Internal/models"
	"github.com/bharatsabne/bookings/Internal/render"
)

// Repo will be used by handlers
var Repo *Repositoy

type Repositoy struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repositoy {
	return &Repositoy{
		App: a,
	}
}

// NewHandler Set the Repository for Handlers
func NewHandler(r *Repositoy) {
	Repo = r
}

func (m *Repositoy) Home(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		//return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)

		fmt.Fprintf(w, "userip: %q is not IP:port", r.RemoteAddr)
	}

	fmt.Println("IP Address: " + ip)
	m.App.Session.Put(r.Context(), "remote_ip", ip)
	render.RenderTemplate(w, r, "Home.page.html", &models.TempateData{})
}

// About Page handler
func (m *Repositoy) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hi, There"
	lol := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = lol
	render.RenderTemplate(w, r, "About.page.html", &models.TempateData{
		StringMap: stringMap,
	})
}

// Contact
func (m *Repositoy) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "Contact.page.html", &models.TempateData{})
}

// Reservation
func (m *Repositoy) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, r, "Make-Reservation.page.html", &models.TempateData{
		Forms: forms.New(nil),
		Data:  data,
	})
}

// Post Reservation hanldes the posting of reservation form
func (m *Repositoy) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("%s", err.Error())
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("FirstName"),
		LastName:  r.Form.Get("LastName"),
		Email:     r.Form.Get("Email"),
		Phone:     r.Form.Get("Phone"),
	}

	form := forms.New(r.PostForm)
	// form.Has("FirstName", r)
	form.Required("FirstName", "LastName", "Email")
	form.MinimumLength("FirstName", 3, r)
	form.MinimumLength("LastName", 3, r)
	form.IsEmail("Email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "Make-Reservation.page.html", &models.TempateData{
			Forms: form,
			Data:  data,
		})
	} else {
		return
	}
}

// Generals
func (m *Repositoy) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "Generals.page.html", &models.TempateData{})
}

// Marjors
func (m *Repositoy) Marjors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "Majors.page.html", &models.TempateData{})
}

// SearchAvailability
func (m *Repositoy) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "Search-Availability.page.html", &models.TempateData{})
}

// POST SearchAvailability
func (m *Repositoy) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start_date")
	end := r.Form.Get("end_date")
	w.Write([]byte(fmt.Sprintf("Start Date is %s and End Date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handles requests for avalibility and send JSON Response
func (m *Repositoy) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	rest := jsonResponse{
		OK:      false,
		Message: "Not Avalible",
	}
	out, err := json.MarshalIndent(rest, "", "     ")
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
