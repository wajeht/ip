package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/oschwald/geoip2-golang"
	"github.com/wajeht/ip/assets"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	f, err := assets.EmbeddedFiles.Open("static/favicon.ico")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "image/x-icon")
	if _, err := io.Copy(w, f); err != nil {
		log.Printf("Error writing favicon: %v", err)
	}
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	f, err := assets.EmbeddedFiles.Open("static/robots.txt")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/plain")
	if _, err := io.Copy(w, f); err != nil {
		log.Printf("Error writing robots.txt: %v", err)
	}
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	json := r.URL.Query().Get("json") == "true" ||
		r.Header.Get("Content-Type") == "application/json"

	if json {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"message": "ok"}`)); err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Ok</title></head><body><span>Ok</span></body></html>")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
	notFound := r.URL.Path != "/"

	geo := r.URL.Query().Get("geo") == "true"

	userAgent := r.Header.Get("User-Agent")

	json := r.URL.Query().Get("format") == "json" ||
		r.URL.Query().Get("json") == "true" ||
		r.Header.Get("Content-Type") == "application/json" ||
		strings.Contains(userAgent, "curl")

	if notFound && json {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte(`{"message": "Not found"}`)); err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	if notFound {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Not found</title></head><body><span>Not found</span></body></html>")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

	ip := getIPAddress(r)

	record, err := LookupLocation(ip)

	if err != nil {
		log.Printf("Error looking up location: %v", err)

		if json {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(`{"error": "Failed to lookup location"}`)); err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("<!DOCTYPE html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>Error</title></head><body><span>Failed to lookup location</span></body></html>")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return
	}

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
		if _, err := w.Write([]byte(response)); err != nil {
			log.Printf("Error writing response: %v", err)
		}
		return

	case json:
		response := fmt.Sprintf(`{"ip": "%s"}`, ip)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(response)); err != nil {
			log.Printf("Error writing response: %v", err)
		}
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
		if _, err := w.Write([]byte(response)); err != nil {
			log.Printf("Error writing response: %v", err)
		}
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
	if _, err := w.Write([]byte(response)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
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

func LookupLocation(ipStr string) (*geoip2.City, error) {
	db, err := geoip2.Open("GeoLite2-City.mmdb")

	if err != nil {
		return nil, err
	}

	defer db.Close()

	ip := net.ParseIP(ipStr)

	if ip == nil {
		return nil, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	record, err := db.City(ip)

	if err != nil {
		return nil, err
	}

	return record, nil
}

const (
	port            = 80
	shutdownTimeout = 30 * time.Second
	readTimeout     = 15 * time.Second
	writeTimeout    = 15 * time.Second
	idleTimeout     = 60 * time.Second
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServer(http.FS(assets.EmbeddedFiles)))

	mux.HandleFunc("GET /favicon.ico", faviconHandler)

	mux.HandleFunc("GET /robots.txt", robotsHandler)

	mux.HandleFunc("GET /healthz", healthzHandler)

	mux.HandleFunc("GET /", ipHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	go func() {
		log.Printf("Server started on port %d", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
