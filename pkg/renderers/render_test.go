package renderers

import (
	"myapp2/pkg/models"
	"net/http"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("Flash value failed")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww myWriter

	err = Render("/home.page.tmpl", &ww, &models.TemplateData{}, r)

	if err != nil {
		t.Error("browser not able to write template", err)
	}

	err = Render("/non_existing.page.tmpl", &ww, &models.TemplateData{}, r)

	if err == nil {
		t.Error("rendering template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r, nil
}
