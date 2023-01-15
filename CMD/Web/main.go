package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bharatsabne/bookings/pkg/config"
	handler "github.com/bharatsabne/bookings/pkg/handlers"
	"github.com/bharatsabne/bookings/pkg/render"
)

// pontNumber is constant
const pontNumber = ":8080"

var app config.AppConfig

var sesson *scs.SessionManager

// main function
func main() {

	//for deply in production set to true
	app.InProduction = false

	sesson = scs.New()
	sesson.Lifetime = 24 * time.Hour
	sesson.Cookie.Persist = true
	sesson.Cookie.SameSite = http.SameSiteLaxMode
	sesson.Cookie.Secure = app.InProduction

	app.Session = sesson

	tc, err := render.CreateTempateCache()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		log.Fatal("Unable to load Template Cahce")
	}
	app.TemplateCache = tc
	app.Usedcache = true

	render.NewTemplates(&app)

	repo := handler.NewRepo(&app)

	handler.NewHandler(repo)

	fmt.Println("Starting application on port ", pontNumber)
	srv := &http.Server{
		Addr:    pontNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}
