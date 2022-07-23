package main

import (
	"html/template"
	"os"
	"path/filepath"

	"github.com/alekslesik/snippetbox.learn/pkg/models"
)

type templateData struct {
    Snippet  *models.Snippet
    Snippets []*models.Snippet
}

func newTemplateCache(dir string) (map[string]*template.Template, error) {

    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exePath := filepath.Dir(ex)
    // fmt.Println(exePath)

    // init new map keeping cache
    cache := map[string]*template.Template{}

    // use func Glob to get all filepathes slice with '.page.html' ext

    pages, err := filepath.Glob(filepath.Join(exePath + dir, "*.page.html"))
    if err != nil {
        return nil, err
    }

    for _, page := range pages {
        // get filename from filepath
        name := filepath.Base(page)

        // handle page
        ts, err := template.ParseFiles(page)
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