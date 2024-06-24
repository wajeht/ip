package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const PORT = 80

	mux := http.NewServeMux()

	// TODO: put this in file server
	mux.HandleFunc("GET /favicon.ico", faviconHandler)

	// TODO: put this in file server
	mux.HandleFunc("GET /robots.txt", robotsHandler)

	mux.HandleFunc("GET /healthz", healthzHandler)

	mux.HandleFunc("GET /", ipHandler)

	log.Println("Server was started on port:", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)

	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
