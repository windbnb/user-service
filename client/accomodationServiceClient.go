package client

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/windbnb/user-service/util"
)

func DeleteAccomodationForHost(hostId uint) error {
	client := &http.Client{}
	accomodationUrl, _ := util.GetAccommodationServicePathRoundRobin()
	req, _ := http.NewRequest("DELETE", accomodationUrl.Next().Host+"/api/accomodation/delete-all/"+strconv.FormatUint(uint64(hostId), 10), nil)

	_, err := client.Do(req)
	if err != nil {
		return errors.New("accomodation service unreachable")
	}

	return nil
}
