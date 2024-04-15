package steamgifts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type GivewayEnteredResponse struct {
	RespType   string `json:"type"`
	EntryCount string `json:"entry_count"`
	Points     string `json:"points"`
}

func (r GivewayEnteredResponse) PointsInt() (int, error) {
	p, err := strconv.Atoi(r.Points)
	if err != nil {
		return -1, errors.New(fmt.Sprintf("Couldn't parse giveaway entered points response: %s", err.Error()))
	}
	return p, nil
}

var emptyResp = GivewayEnteredResponse{}

func Enter(details GivewayDetailsPage) (GivewayEnteredResponse, error) {
	client := &http.Client{
		Jar:       newCookieJar(),
		Transport: &http.Transport{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	formData := url.Values{}
	formData.Add("do", "entry_insert")
	formData.Add("xsrf_token", details.XsrfToken)
	formData.Add("code", details.Code)

	req, err := http.NewRequest("POST", steamGiftsUrl+"/ajax.php", strings.NewReader(formData.Encode()))
	if err != nil {
		return GivewayEnteredResponse{}, nil
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Println("Sending request...")
	resp, err := client.Do(req)
	log.Println("Done.")
	if err != nil {
		return emptyResp, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var enteredResp GivewayEnteredResponse
	json.Unmarshal(bodyBytes, &enteredResp)
	if err != nil {
		return emptyResp, err
	}
	if enteredResp.RespType == "" {
		return enteredResp, notLoggedIn
	}
	return enteredResp, nil
}
