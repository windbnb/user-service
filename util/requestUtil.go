package util

import (
	"net/url"

	roundrobin "github.com/hlts2/round-robin"
)

var BaseAccomodationServicePathRoundRobin, _ = roundrobin.New(
    &url.URL{Host: "http://localhost:8082/api"},
)

var BaseReservationServicePathRoundRobin, _ = roundrobin.New(
    &url.URL{Host: "http://localhost:8083/api"},
)
