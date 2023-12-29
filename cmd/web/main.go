package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	// define command line flags
	addr := flag.String("addr", ":4000", "HTTP network address")

	// parse all command line flags
	flag.Parse()

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

	log.Printf("starting server on http://localhost%s", *addr)

	// use http.ListenAndServe() func to start web server.
	// pass in address and servemux, log error and exit.
	err := http.ListenAndServe(*addr, mux)
	log.Fatal(err)
}
