package util

import (
	"net/url"
	"os"

	roundrobin "github.com/hlts2/round-robin"
)

func GetReservationServicePathRoundRobin() (roundrobin.RoundRobin, error) {
	reservationServicePath, reservationServicePathFound := os.LookupEnv("RESERVATION_SERVICE_PATH")
	if !reservationServicePathFound {
		reservationServicePath = "http://localhost:8083"
	}

	return roundrobin.New(
		&url.URL{Host: reservationServicePath},
	)
}

func GetAccommodationServicePathRoundRobin() (roundrobin.RoundRobin, error) {
	accommodationServicePath, accommodationServicePathFound := os.LookupEnv("ACCOMMODATION_SERVICE_PATH")
	if !accommodationServicePathFound {
		accommodationServicePath = "http://localhost:8082"
	}

	return roundrobin.New(
		&url.URL{Host: accommodationServicePath},
	)
}
