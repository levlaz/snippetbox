package main

import (
	"html/template"
	"path/filepath"

	"snippetbox.levlaz.org/internal/models"
)

type templateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
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

		// create slice containing filepaths for our base template, partials, and the page
		files := []string{
			"./ui/html/base.tmpl.html",
			"./ui/html/partials/nav.tmpl.html",
			page,
		}

		// parse files into template set
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		// add template set to the map using name of page as key
		cache[name] = ts
	}

	return cache, nil
}
