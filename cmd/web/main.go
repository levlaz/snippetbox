package main

import (
	"log"
	"net/http"
)

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
