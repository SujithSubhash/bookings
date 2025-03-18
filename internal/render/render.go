package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/SujithSubhash/bookings/internal/config"
	"github.com/SujithSubhash/bookings/internal/models"
	"github.com/justinas/nosurf"
)

var function = template.FuncMap{}
var app *config.AppConfig

// sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

// NewTemplates sets the config for the template package
func AddDefaultData(td *models.TemplateData,r *http.Request) *models.TemplateData {
	td.CSRFToken= nosurf.Token(r)
	return td
}

// render templates using html template
func RenderTemplate(w http.ResponseWriter,r *http.Request, html string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {

		//create a template cache
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get requetsted template from cache
	t, ok := tc[html]
	if !ok {
		log.Fatal("could not get template from template cache")
	}
	buf := new(bytes.Buffer)

	td = AddDefaultData(td,r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("error writing template to browser", err)
	}

	//render the template

	// 	_, err = buf.WriteTo(w)
	// 	if err != nil {
	// 		log.Println(err)

}
//createTemplateCache create a template cache as a map 
func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	//get all of the files named *.page.html from ./templates

	pages, err := filepath.Glob("./templates/*.page.html")

	if err != nil {
		return myCache, err
	}
	//range through all the files ending with page.html
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(function).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts
	}
	return myCache, nil
}
