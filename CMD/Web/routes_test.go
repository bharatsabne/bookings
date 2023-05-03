package main

import (
	"fmt"
	"testing"

	"github.com/bharatsabne/bookings/Internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		var str = fmt.Sprintf("Did not receive http.handler, insted received %T", v)
		t.Error(str)
	}
}
