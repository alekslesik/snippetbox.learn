package main

import (
	"html/template"
	"path/filepath"
	"time"

	"github.com/alekslesik/snippetbox.learn/pkg/models"
)

type templateData struct {
	CurrentYear int
	Snippet     *models.Snippet
	Snippets    []*models.Snippet
}

// Return nicely formatted string of time.Time object
func humanDate(t time.Time) string {
	return t.Format("02 Jan 2006 at 15:04")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

	// init new map keeping cache
	cache := map[string]*template.Template{}

	// use func Glob to get all filepathes slice with '.page.html' ext

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// get filename from filepath
		name := filepath.Base(page)

		// The template.FuncMap must be registered with the template set before
		// call the ParseFiles() method. This means we have to use template.New
		// create an empty template set, use the Funcs() method t
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// use ParseGlob to add all frame patterns (base.layout.html)
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// add received patterns set to cache, using page name
		// (ext home.page.html) as a key for our map
		cache[name] = ts
	}

	return cache, nil

}
