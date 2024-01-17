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
		http.NotFound(w, r)
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
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute() on template set to write the template as response
	// body. Last param is dynamic data.
	err = ts.ExecuteTemplate(w, "base", nil)
	if err != nil {
		app.logger.Error(err.Error(), "method", r.Method, "uri", r.URL.RequestURI())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// get value of id paramter from query string, try to convert
	// to int. If cannot convert or value is less than 1, return 404
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// use fmt.Fprintf() to interpolate id value with response
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}

func (app *application) snippetCreate(w http.ResponseWriter, r *http.Request) {
	// use r.Method to check if request is using POST
	if r.Method != "POST" {
		// let client know which methods are allowed
		w.Header().Set("Allow", http.MethodPost)

		// use http.Error() function to send 405 and
		// method not allowed
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}
