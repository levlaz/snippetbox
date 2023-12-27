package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// create file server to serve files out of "./ui/static" dir
	// note path is relative to project dir root
	fileServer := http.FileServer(http.Dir("./ui/static"))

	// register file server as handler for all URL paths that
	// start with "static"
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// register other application routes
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Print("starting server on http://localhost:4000")

	// use http.ListenAndServe() func to start web server.
	// pass in address and servemux, log error and exit.
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
