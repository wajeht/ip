package main

import (
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

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
