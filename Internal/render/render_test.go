package render

import (
	"testing"

	"github.com/bharatsabne/bookings/Internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td *models.TempateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	testSession.Put(r.Context(), "flash", "This is flash Message")
	result := AddDefaultData(td, r)
	if result.Flash == "Test" {
		t.Error("Failed")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplate = "./../../Templates"
	tc, err := CreateTempateCache()
	if err != nil {
		t.Error(err)
	}
	app.TemplateCache = tc

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	var ww myWriter

	err = RenderTemplate(&ww, r, "Home.page.html", &models.TempateData{})
	if err != nil {
		t.Error(err)
	}
	err = RenderTemplate(&ww, r, "NotExisted.page.html", &models.TempateData{})
	if err == nil {
		t.Fatal("Rendered page that does not existed")
	}
}

func TestNewTemplate(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCatch(t *testing.T) {
	pathToTemplate = "./../../Templates"
	_, err := CreateTempateCache()
	if err != nil {
		t.Error(err)
	}
}
