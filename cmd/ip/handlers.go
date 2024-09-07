package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

// TODO: put this in file server
func robotsHandler(w http.ResponseWriter, r *http.Request) {
	basePath, _ := os.Getwd()
	filePath := filepath.Join(basePath, "web/static/robots.txt")
	http.ServeFile(w, r, filePath)
}

// TODO: put this in file server
func faviconHandler(w http.ResponseWriter, r *http.Request) {
	basePath, _ := os.Getwd()
	filePath := filepath.Join(basePath, "web/static/favicon.ico")
	http.ServeFile(w, r, filePath)
}

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
		formattedGeo := fmt.Sprintf(`<strong>IP:</strong> %s<br>`, ip)
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
					<script defer data-domain="ip.jaw.dev" src="https://plausible.jaw.dev/js/script.js"></script>
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

	response := fmt.Sprintf(`
		<!DOCTYPE html>
		<html lang="en">
		<head>
				<meta charset="UTF-8">
				<script defer data-domain="ip.jaw.dev" src="https://plausible.jaw.dev/js/script.js"></script>
				<title>ip</title>
		</head>
		<body>
				<pre>%s</pre>
		</body>
		</html>`, ip)
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(response))
}
