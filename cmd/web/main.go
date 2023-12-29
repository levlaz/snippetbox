package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// TODO: make this configurable
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

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

	logger.Info("starting server", slog.String("addr", ":4000"))

	// use http.ListenAndServe() func to start web server.
	// pass in address and servemux, log error and exit.
	err := http.ListenAndServe(*addr, mux)
	logger.Error(err.Error())
	os.Exit(1)
}
