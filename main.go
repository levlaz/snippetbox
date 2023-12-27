package main

import (
	"log"
	"net/http"
)

// define home handler function which writes
// a byte slice containing hello from slopshop
// as response body
func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from slopshop"))
}

// Add snippetView handler function
func snippetView(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("display a specific snippet..."))
}

// Add a snippetCreate handler
func snippetCreate(w http.ResponseWriter, r *http.Request) {
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
