package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

const PORT = 80

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	isJSON := r.URL.Query().Get("json") == "true" ||
		r.Header.Get("Content-Type") == "application/json"

	if isJSON {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<html><body><span>ok</span></body></html>"))
}

func getIPAddress(r *http.Request) string {
	ipAddress := r.Header.Get("x-forwarded-for")
	if ipAddress == "" {
		ipAddress, _, _ = net.SplitHostPort(r.RemoteAddr)
	} else {
		ips := strings.Split(ipAddress, ", ")
		ipAddress = ips[0]
	}
	return ipAddress
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIPAddress(r)

	geo := r.URL.Query().Get("geo") == "true"
	json := r.URL.Query().Get("json") == "true" ||
		r.URL.Query().Get("format") == "json" ||
		r.Header.Get("Content-Type") == "application/json"

	if json && geo {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"ip": "%s"}`, ip)))
		return
	}

	if json {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf(`{"ip": "%s"}`, ip)))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><span>%s</span></body></html>", ip)))
}

func main() {
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/", ipHandler)
	fmt.Println("Server was started on http://localhost:", PORT)
	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
