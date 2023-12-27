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

func main() {
	// use http.NewServeMux() func to init a new servemux, then
	// register home function as handler for "/" URL pattern
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	// print log message to say the server is up
	log.Print("starting server on :4000")

	// use http.ListenAndServe() func to start web server.
	// pass in address and servemux, log error and exit.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
