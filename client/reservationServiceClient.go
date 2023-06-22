package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/util"
)

func CheckReservations(userId uint, role string, tokenString string) error {
	if role != "owner" && role != "guest" {
		return errors.New("invalid role specified")
	}

	client := &http.Client{}
	reservationUrl, _ := util.GetReservationServicePathRoundRobin()
	req, _ := http.NewRequest("GET", reservationUrl.Next().Host+"/api/reservationRequest/"+role+"/"+strconv.FormatUint(uint64(userId), 10), nil)

	req.Header.Set("Authorization", tokenString)

	response, err := client.Do(req)
	if err != nil {
		return errors.New("reservation service unreachable")
	}
	defer response.Body.Close()

	var reservations []model.ReservationRequestDto
	err = json.NewDecoder(response.Body).Decode(&reservations)
	if err != nil {
		return errors.New("failed to parse active reservation list")
	}

	if len(reservations) > 0 {
		return errors.New(role + " has active reservations")
	}

	return nil
}
