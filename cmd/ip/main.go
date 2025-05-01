package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wajeht/ip/assets"
)

const PORT = 80

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(assets.EmbeddedFiles))

	mux.Handle("GET /static/", fileServer)

	mux.HandleFunc("GET /favicon.ico", faviconHandler)

	mux.HandleFunc("GET /robots.txt", robotsHandler)

	mux.HandleFunc("GET /healthz", healthzHandler)

	mux.HandleFunc("GET /", ipHandler)

	log.Println("Server was started on port:", PORT)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux))
}
