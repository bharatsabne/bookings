package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"generals-quarters", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"marjors-suite", "/marjors-suite", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"reservation-summary", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"search-availability", "/search-availability", "POST", []postData{
		{key: "start_date", value: "2023-05-04"},
		{key: "end_date", value: "2023-05-08"},
	}, http.StatusOK},
	{"search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start_date", value: "2023-05-04"},
		{key: "end_date", value: "2023-05-08"},
	}, http.StatusOK},
	{"make-reservation", "/make-reservation", "POST", []postData{
		{key: "FirstName", value: "Rudransh"},
		{key: "LastName", value: "Sabne"},
		{key: "Email", value: "rudransh.sabne@website.com"},
		{key: "Phone", value: "5555-4444-3333"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d ", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else {
			values := url.Values{}

			for _, x := range e.params {
				values.Add(x.key, x.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d ", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		}
	}
}
