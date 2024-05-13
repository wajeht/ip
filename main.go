package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/oschwald/geoip2-golang"
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

func LookupLocation(ipStr string) *geoip2.City {
	db, err := geoip2.Open("GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ip := net.ParseIP(ipStr)

	if ip == nil {
		log.Fatalf("Invalid IP address: %s", ipStr)
	}

	record, err := db.City(ip)
	if err != nil {
		log.Fatal(err)
	}

	return record
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	ip := getIPAddress(r)

	geo := r.URL.Query().Get("geo") == "true"
	json := r.URL.Query().Get("json") == "true" ||
		r.URL.Query().Get("format") == "json" ||
		r.Header.Get("Content-Type") == "application/json"

	if json && geo {
		record := LookupLocation(ip)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`{
        "ip": "%s",
        "range": [%f, %f],
        "country": "%s",
        "region": "%s",
        "eu": %t,
        "timezone": "%s",
        "city": "%s",
        "ll": [%f, %f],
        "metro": %d,
        "area": %d
    }`, ip, record.Location.Longitude, record.Location.Longitude,
			record.Country.IsoCode, record.Subdivisions[0].IsoCode,
			record.Country.IsInEuropeanUnion, record.Location.TimeZone,
			record.City.Names["en"], record.Location.Latitude, record.Location.Longitude,
			record.Location.MetroCode, record.Location.AccuracyRadius)
		w.Write([]byte(response))
		return
	}

	if json {
		record := LookupLocation(ip)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := fmt.Sprintf(`{
        "ip": "%s",
        "range": [%f, %f],
        "country": "%s",
        "region": "%s",
        "eu": %t,
        "timezone": "%s",
        "city": "%s",
        "ll": [%f, %f],
        "metro": %d,
        "area": %d
    }`, ip, record.Location.Longitude, record.Location.Longitude,
			record.Country.IsoCode, record.Subdivisions[0].IsoCode,
			record.Country.IsInEuropeanUnion, record.Location.TimeZone,
			record.City.Names["en"], record.Location.Latitude, record.Location.Longitude,
			record.Location.MetroCode, record.Location.AccuracyRadius)
		w.Write([]byte(response))
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
