// Home Page Handler
package handler

import (
	"fmt"
	"net"
	"net/http"

	"github.com/bharatsabne/bookings/pkg/config"
	"github.com/bharatsabne/bookings/pkg/models"
	"github.com/bharatsabne/bookings/pkg/render"
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
	render.RenderTemplate(w, "Home.page.html", &models.TempateData{})
}

// About Page handler
func (m *Repositoy) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hi, There"
	lol := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = lol
	render.RenderTemplate(w, "About.page.html", &models.TempateData{
		StringMap: stringMap,
	})
}
