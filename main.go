package main

import (
	"log"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	// check if current request URL path exactly matches "/".
	// if not, use http.NotFound() function to send 404.
	// we need to do this because my default servemux treates
	// subtree paths as catch-alls
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Write([]byte("Hello from slopshop"))
}

// Add snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a specific snippet..."))
}

// Add a snippetCreate handler
func snippetCreate(w http.ResponseWriter, r *http.Request) {
	// use r.Method to check if request is using POST
	if r.Method != "POST" {
		// if not, send 405 status, method not allowed
		// we return so the subsequent code is not executed
		// use Header.Set() to add Alllow:POST header to response
		w.Header().Set("Allow", "POST")
		w.WriteHeader(405)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func main() {
	// use http.NewServeMux() func to init a new servemux, then
	// register home function as handler for "/" URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	// print log message to say the server is up
	log.Print("starting server on http://localhost:4000")

	// use http.ListenAndServe() func to start web server.
	// pass in address and servemux, log error and exit.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
