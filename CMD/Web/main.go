package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bharatsabne/bookings/Internal/config"
	handler "github.com/bharatsabne/bookings/Internal/handlers"
	"github.com/bharatsabne/bookings/Internal/models"
	"github.com/bharatsabne/bookings/Internal/render"
)

// pontNumber is constant
const pontNumber = ":8080"

var app config.AppConfig

var session *scs.SessionManager

// main function
func main() {

	err := run()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting application on port ", pontNumber)
	srv := &http.Server{
		Addr:    pontNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
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

	tc, err := render.CreateTempateCache()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		log.Fatal("Unable to load Template Cahce")
		return err
	}
	app.TemplateCache = tc
	app.Usedcache = false

	repo := handler.NewRepo(&app)
	handler.NewHandler(repo)
	render.NewTemplates(&app)
	return nil
}
