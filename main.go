package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

const PORT = 80

func healthzHandler(w http.ResponseWriter, r *http.Request) {
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
	ipAddress := getIPAddress(r)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><span>%s</span></body></html>", ipAddress)))
}

func main() {
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/", ipHandler)
	fmt.Println("Server was started on http://localhost:", PORT)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
