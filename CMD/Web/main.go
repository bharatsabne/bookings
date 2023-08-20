package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	helpers "github.com/bharatsabne/bookings/Internal/Helpers"
	"github.com/bharatsabne/bookings/Internal/config"
	driver "github.com/bharatsabne/bookings/Internal/drivers"
	handler "github.com/bharatsabne/bookings/Internal/handlers"
	"github.com/bharatsabne/bookings/Internal/models"
	"github.com/bharatsabne/bookings/Internal/render"
)

// pontNumber is constant
const pontNumber = ":8081"

var app config.AppConfig

var session *scs.SessionManager

var infoLog *log.Logger
var errorLog *log.Logger

// main function
func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}
	defer db.SQL.Close()

	fmt.Println("Starting application on port ", pontNumber)
	srv := &http.Server{
		Addr:    pontNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	//what am I going to put in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	//for deply in production set to true
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	//connect to database
	log.Println("Connecting to database")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=Pass@1234")
	if err != nil {
		log.Fatal("Con't connect to database")
	}
	log.Printf("Connected to database")
	tc, err := render.CreateTempateCache()
	if err != nil {
		fmt.Println(fmt.Sprint(err))
		log.Fatal("Unable to load Template Cahce")
		return nil, err
	}
	app.TemplateCache = tc
	app.Usedcache = false

	repo := handler.NewRepo(&app, db)
	handler.NewHandler(repo)
	render.NewRederer(&app)
	helpers.NewHelper(&app)
	return db, nil
}
