package main

import (
	"fmt"
	"log"
	"net/http"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./favicon.ico")
}

func main() {
	const PORT = 80

	mux := http.NewServeMux()

	mux.HandleFunc("GET /favicon.ico", faviconHandler)

	mux.HandleFunc("GET /healthz", healthzHandler)

	mux.HandleFunc("GET /", ipHandler)

	log.Println("Server was started on port:", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)

	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}