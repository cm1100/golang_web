package renderers

import (
	"fmt"
	"log"
	"myapp2/pkg/config"
	"myapp2/pkg/models"
	"net/http"
	"path/filepath"
	"time"

	//"fmt"
	"html/template"

	"github.com/justinas/nosurf"
	//"reflect"
)

var app *config.AppConfig

var functions = template.FuncMap{
	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

var pathToTemplates = "./templates"

func SetAppConfig(a *config.AppConfig) {
	app = a
}

func HumanDate(t time.Time) string {

	return t.Format("2006-01-02")

}

func FormatDate(t time.Time, f string) string {

	return t.Format(f)
}

// Return slice of int from 1 going to a number
func Iterate(count int) []int {

	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

func Add(a, b int) int {
	return a + b
}

func SetPath(s string) {
	pathToTemplates = s
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.CSRFToken = nosurf.Token(r)
	if app.Session.Exists(r.Context(), "user_id") {
		td.IsAuthenticated = 1
	}
	return td
}

func Render(tmpl string, w http.ResponseWriter, temp_data *models.TemplateData, r *http.Request) error {

	temp_data = AddDefaultData(temp_data, r)
	t, err := template.ParseFiles(pathToTemplates+tmpl, fmt.Sprintf("%s/base.layout.tmpl", pathToTemplates))
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, filepath.Base(tmpl), temp_data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
	//fmt.Println(reflect.TypeOf(t))
}

func RenderAdmin(tmpl string, w http.ResponseWriter, temp_data *models.TemplateData, r *http.Request) error {

	temp_data = AddDefaultData(temp_data, r)
	name := filepath.Base(tmpl)
	t, err := template.New(name).Funcs(functions).ParseFiles(pathToTemplates+tmpl, fmt.Sprintf("%s/admin_base.layout.tmpl", pathToTemplates))
	if err != nil {
		log.Println(err)
		return err
	}

	err = t.ExecuteTemplate(w, filepath.Base(tmpl), temp_data)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
	//fmt.Println(reflect.TypeOf(t))
}
