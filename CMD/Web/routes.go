package main

import (
	"net/http"

	"github.com/bharatsabne/bookings/Internal/config"
	handler "github.com/bharatsabne/bookings/Internal/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	// mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)
	mux.Get("/", handler.Repo.Home)
	mux.Get("/about", handler.Repo.About)
	mux.Get("/contact", handler.Repo.Contact)
	mux.Get("/generals-quarters", handler.Repo.Generals)
	mux.Get("/marjors-suite", handler.Repo.Marjors)

	mux.Get("/search-availability", handler.Repo.SearchAvailability)
	mux.Post("/search-availability", handler.Repo.PostSearchAvailability)
	mux.Post("/search-availability-json", handler.Repo.AvailabilityJSON)

	mux.Get("/make-reservation", handler.Repo.Reservation)
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
