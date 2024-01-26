package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.levlaz.org/internal/models"
)

type templateData struct {
	CurrentYear int
	Snippet     models.Snippet
	Snippets    []models.Snippet
}

func newTemplateCache() (map[string]*template.Template, error) {
	// init a new map to act as cache
	cache := map[string]*template.Template{}

	// use filepath.Glob() to get slice of all filepaths that match
	// pattern. Give us slice of all filepaths for applicaiton page
	// templates.
	pages, err := filepath.Glob("./ui/html/pages/*.tmpl.html")
	if err != nil {
		return nil, err
	}

	// loop through filepaths
	for _, page := range pages {
		// extract filename
		name := filepath.Base(page)

		// parse base template file into template set
		ts, err := template.ParseFiles("./ui/html/base.tmpl.html")
		if err != nil {
			return nil, err
		}

		// call ParseGlob() on this template set to add any partials
		ts, err = ts.ParseGlob("./ui/html/partials/*.tmpl.html")
		if err != nil {
			return nil, err
		}

		// call ParseFiles() on this template set to add the page template
		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add template set to the map using name of page as key
		cache[name] = ts
	}

	return cache, nil
}
