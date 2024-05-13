package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	json := r.URL.Query().Get("json") == "true" ||
		r.Header.Get("Content-Type") == "application/json"

	if json {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "ok"}`))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Ok</title></head><body><span>Ok</span></body></html>"))
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
	notFound := r.URL.Path != "/"

	geo := r.URL.Query().Get("geo") == "true"

	json := r.URL.Query().Get("format") == "json" ||
		r.URL.Query().Get("json") == "true" ||
		r.Header.Get("Content-Type") == "application/json"

	if notFound && json {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{ message: "Not found" }`))
		return
	}

	if notFound {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Not found</title></head><body><span>Not found</span></body></html>"))
		return
	}

	ip := getIPAddress(r)

	record := LookupLocation(ip)

	switch {
	case json && geo:
		response := fmt.Sprintf(`{"ip": "%s", "range": [%f, %f], "country": "%s", "region": "%s", "eu": %t, "timezone": "%s", "city": "%s", "ll": [%f, %f], "metro": %d, "area": %d}`,
			ip, record.Location.Longitude, record.Location.Latitude,
			record.Country.IsoCode, record.Subdivisions[0].IsoCode,
			record.Country.IsInEuropeanUnion, record.Location.TimeZone,
			record.City.Names["en"], record.Location.Longitude, record.Location.Latitude,
			record.Location.MetroCode, record.Location.AccuracyRadius)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return

	case json:
		response := fmt.Sprintf(`{"ip": "%s"}`, ip)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return

	case geo:
		formattedGeo := fmt.Sprintf(`<strong>ip:</strong> %s<br>`, ip)
		formattedGeo += fmt.Sprintf(`<strong>Country:</strong> %s<br>`, record.Country.IsoCode)
		formattedGeo += fmt.Sprintf(`<strong>Region:</strong> %s<br>`, record.Subdivisions[0].IsoCode)
		formattedGeo += fmt.Sprintf(`<strong>City:</strong> %s<br>`, record.City.Names["en"])
		formattedGeo += fmt.Sprintf(`<strong>Latitude:</strong> %f<br>`, record.Location.Latitude)
		formattedGeo += fmt.Sprintf(`<strong>Longitude:</strong> %f<br>`, record.Location.Longitude)
		formattedGeo += fmt.Sprintf(`<strong>Timezone:</strong> %s<br>`, record.Location.TimeZone)
		formattedGeo += fmt.Sprintf(`<strong>Metro Code:</strong> %d<br>`, record.Location.MetroCode)
		formattedGeo += fmt.Sprintf(`<strong>Area Code:</strong> %d<br>`, record.Location.AccuracyRadius)
		response := fmt.Sprintf(`
			<!DOCTYPE html>
			<html lang="en">
			<head>
					<meta charset="UTF-8">
					<title>IP</title>
			</head>
			<body>
					<pre>%s</pre>
			</body>
			</html>`, formattedGeo)
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ip))
}

func main() {
	const PORT = 8080

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthzHandler)

	mux.HandleFunc("GET /", ipHandler)

	log.Println("Server was started on port:", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), mux)

	if err != nil {
		log.Fatalf("Error starting server: %s", err)
	}
}
