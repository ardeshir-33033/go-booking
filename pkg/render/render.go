package render

import (
	"bytes"
	"fmt"
	"github.com/ardeshir-33033/go-booking/pkg/config"
	"github.com/ardeshir-33033/go-booking/pkg/models"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

var functions = template.FuncMap{}

var app *config.AppConfig

func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template
	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	ts, ok := tc[tmpl]
	if !ok {
		log.Fatal("Cannot retrieve template cache")
		return
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)
	_ = ts.Execute(buf, td)
	_, err := buf.WriteTo(w)
	if err != nil {
		log.Fatal(w, err)
		return
	}

	parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
	err = parsedTemplate.Execute(w, nil)

	if err != nil {
		fmt.Println("error ", err)
		return
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// public pages
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// Loop through the pages one-by-one.
	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println(err)
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			fmt.Println(err)
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				fmt.Println(err)
				return myCache, err
			}
		}

		// Add the template set to the cache,
		myCache[name] = ts
	}

	return myCache, nil
}
