package main

import (
	"fmt"
	"net/http"
)

const PORT = 80

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><span>ok</span></body></html>"))
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/", ipHandler)
	fmt.Println("Server was started on http://localhost:", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
