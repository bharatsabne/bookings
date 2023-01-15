package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

//	func WriteToConsole(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			fmt.Println("Hits the page")
//			next.ServeHTTP(w, r)
//		})
//	}
//
// NoSurf adds CSFR protection to all POST request
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
