package models

import "github.com/bharatsabne/bookings/Internal/forms"

// TempateData holds data send fromhandlers to templates
type TempateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Forms     *forms.Form
}
