package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/bharatsabne/bookings/Internal/config"
)

var app *config.AppConfig

// setup app config for helpers
func NewHelper(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client Error with status of ", status)
	http.Error(w, http.StatusText(status), status)
}
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
