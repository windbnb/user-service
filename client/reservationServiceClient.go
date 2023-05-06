package client

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/windbnb/user-service/model"
	"github.com/windbnb/user-service/util"
)

func CheckReservations(userId uint, role string) error {
	if role != "owner" && role != "guest" {
		return errors.New("invalid role specified")
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", util.BaseReservationServicePathRoundRobin.Next().Host+"/reservationRequest/" + role + "/"+strconv.FormatUint(uint64(userId), 10), nil)

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