package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bharatsabne/bookings/Internal/config"
	"github.com/bharatsabne/bookings/Internal/models"
)

var testSession *scs.SessionManager
var testApp config.AppConfig

// func TestMain(m *testing.M) {
// 	//what am I going to put in the session
// 	gob.Register(models.Reservation{})
// 	//for deply in production set to true
// 	testApp.InProduction = false

// 	testSession = scs.New()
// 	testSession.Lifetime = 24 * time.Hour
// 	testSession.Cookie.Persist = true
// 	testSession.Cookie.SameSite = http.SameSiteLaxMode
// 	testSession.Cookie.Secure = false

//		testApp.Session = testSession
//		app = &testApp
//	}
func TestMain(m *testing.M) {
	//what am I going to put in the session
	gob.Register(models.Reservation{})
	//for deply in production set to true
	testApp.InProduction = false
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "Error\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog
	testSession = scs.New()
	testSession.Lifetime = 24 * time.Hour
	testSession.Cookie.Persist = true
	testSession.Cookie.SameSite = http.SameSiteLaxMode
	testSession.Cookie.Secure = false
	testApp.Session = testSession
	app = &testApp
}
func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		return nil, err
	}
	ctx := r.Context()
	ctx, _ = testSession.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}

type myWriter struct{}

func (tw *myWriter) Header() http.Header {
	var v http.Header
	return v
}

func (tw *myWriter) WriteHeader(statusCode int) {
}
func (tw *myWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}
