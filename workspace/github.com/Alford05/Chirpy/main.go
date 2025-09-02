package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	fileHandler := http.FileServer(http.Dir("."))
	mux.Handle("/", fileHandler)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Starting server on port %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
