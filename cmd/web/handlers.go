package main

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	// check if current request URL path exactly matches "/".
	// if not, use http.NotFound() function to send 404.
	// we need to do this because my default servemux treates
	// subtree paths as catch-alls
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	// init slice containing paths to template files
	files := []string{
		"./ui/html/base.tmpl.html",
		"./ui/html/partials/nav.tmpl.html",
		"./ui/html/pages/home.tmpl.html",
	}

	// use template.ParseFiles() function to read template file
	// or return 500 error
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		app.serverError(w, r, err)
		return
	}

	// Execute() on template set to write the template as response
	// body. Last param is dynamic data.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		app.serverError(w, r, err)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// get value of id paramter from query string, try to convert
	// to int. If cannot convert or value is less than 1, return 404
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		app.logger.Debug(r.URL.RawQuery)
		app.notFound(w)
		return
	}

	// use fmt.Fprintf() to interpolate id value with response
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// use r.Method to check if request is using POST
	if r.Method != http.MethodPost {
		// let client know which methods are allowed
		w.Header().Set("Allow", http.MethodPost)

		// use http.Error() function to send 405 and
		// method not allowed
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	title := "0 snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := 7

	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/snippet/view?id=%d", id), http.StatusSeeOther)
}
