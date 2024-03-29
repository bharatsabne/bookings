package render

import (
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/bharatsabne/bookings/Internal/config"
	"github.com/bharatsabne/bookings/Internal/models"
	"github.com/justinas/nosurf"
)

var app *config.AppConfig

var pathToTemplate = "./Templates"

var functions = template.FuncMap{}

// NewRederer Set the config for the template package
func NewRederer(a *config.AppConfig) {
	app = a
}
func AddDefaultData(td *models.TempateData, r *http.Request) *models.TempateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// Template Renders the template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TempateData) error {
	//get themplate cache from app.config

	//create template cache
	// tc, err := CreateTempateCache()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	var tc map[string]*template.Template
	if app.Usedcache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTempateCache()
	}

	//get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template, from Template Cache")
		return errors.New("Could not get template, from Template Cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)
	err := t.Execute(buf, td)
	if err != nil {
		log.Println(err)
		return err
	}
	//render the templeate
	_, err = buf.WriteTo(w)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
func CreateTempateCache() (map[string]*template.Template, error) {
	// myChach := make(map[string]*template.Template)
	myChach := map[string]*template.Template{}

	//get all templates from ./Template folder which are *.page.html
	pages, err := filepath.Glob(pathToTemplate + "/*.page.html")
	if err != nil {
		return myChach, err
	}
	//range through all files
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myChach, err
		}

		matches, err := filepath.Glob(pathToTemplate + "/*.layout.html")
		if err != nil {
			return myChach, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(pathToTemplate + "/*.layout.html")
			if err != nil {
				return myChach, err
			}
		}
		myChach[name] = ts
	}
	return myChach, nil
}

// func RenderTemplateTest(w http.ResponseWriter, tmpl string) {
// 	parsedTemplate, _ := template.ParseFiles("./Templates/"+tmpl, "./Templates/Base.layout.html")
// 	err := parsedTemplate.Execute(w, nil)
// 	if err != nil {
// 		fmt.Println("Error Parsing template")
// 		return
// 	}
// }

// var tc = make(map[string]*template.Template)

// func RenderTemplate(w http.ResponseWriter, t string) {
// 	var tmpl *template.Template
// 	var err error
// 	//check to see if we have already have the template in our cache
// 	_, inMap := tc[t]
// 	if !inMap {
// 		log.Println("Creating template and adding to cache")
// 		//need to create templaet
// 		err = CreateTemplateCatche(t)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		//already have
// 		log.Println("Using cache template")
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	}
// 	tmpl = tc[t]
// 	err = tmpl.Execute(w, nil)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }
// func CreateTemplateCatche(t string) error {
// 	templates := []string{
// 		fmt.Sprintf("./Templates/%s", t),
// 		"./Templates/Base.layout.html",
// 	}
// 	//parse the template
// 	tmpl, err := template.ParseFiles(templates...)
// 	if err != nil {
// 		return err
// 	}
// 	//add template to cache
// 	tc[t] = tmpl
// 	return nil
// }
