package handler

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bharatsabne/bookings/Internal/config"
	"github.com/bharatsabne/bookings/Internal/models"
	"github.com/bharatsabne/bookings/Internal/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplate = "./../../Templates"
var functions = template.FuncMap{}

func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// SessionLoad loads and saved sessions per request
func SessionLoad(next http.Handler) http.Handler {
	return app.Session.LoadAndSave(next)
}

func getRoutes() http.Handler {
	//what am I going to put in the session
	gob.Register(models.Reservation{})
	//for deply in production set to true
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	tc, err := CreateTestTempateCache()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		log.Fatal("Unable to load Template Cahce")
	}
	app.TemplateCache = tc
	app.Usedcache = true

	repo := NewRepo(&app)
	NewHandler(repo)
	render.NewTemplates(&app)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	// mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/marjors-suite", Repo.Marjors)

	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}

func CreateTestTempateCache() (map[string]*template.Template, error) {
	// myChach := make(map[string]*template.Template)
	myChach := map[string]*template.Template{}

	//get all templates from ./Template folder which are *.page.html
	pages, err := filepath.Glob(pathToTemplate + "/*.page.html")
	if err != nil {
		return myChach, err
	}
	//range through all files
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myChach, err
		}

		matches, err := filepath.Glob(pathToTemplate + "/*.layout.html")
		if err != nil {
			return myChach, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(pathToTemplate + "/*.layout.html")
			if err != nil {
				return myChach, err
			}
		}
		myChach[name] = ts
	}
	return myChach, nil
}
