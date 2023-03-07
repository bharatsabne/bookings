package forms

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Forms Craete a custom form struct, embeds url.values object
type Form struct {
	url.Values
	Errors errors
}

// Valid Returns true if ther are no error otherwise false
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New Iniliaze new from struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has checks if form filed is in post and not empty
func (f *Form) Has(filed string, r *http.Request) bool {
	x := r.Form.Get(filed)
	if x == "" {
		f.Errors.Add(filed, "This field cannot be blank")
		return false
	}
	return true
}

// Required checks for equired fileds
func (f *Form) Required(fileds ...string) {
	for _, filed := range fileds {
		value := f.Get(filed)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(filed, "This filed cannot be blank")
		}
	}
}

// MinimumLength checks minimum length
func (f *Form) MinimumLength(filed string, length int, r *http.Request) bool {
	filedLength := len(r.Form.Get(filed))
	if filedLength < length {
		f.Errors.Add(filed, fmt.Sprintf("This filed must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail check for valid email address
func (f *Form) IsEmail(filed string) {
	if !govalidator.IsEmail(f.Get(filed)) {
		f.Errors.Add(filed, "Invalid Email Address")
	}
}
