package steamgifts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type steamGiftsRequest interface {
	asFormData() *url.Values
}

type GivewayEnteredResponse struct {
	RespType   string `json:"type"`
	EntryCount string `json:"entry_count"`
	Points     string `json:"points"`
}

func (r GivewayEnteredResponse) PointsInt() (int, error) {
	p, err := strconv.Atoi(r.Points)
	if err != nil {
		return -1, fmt.Errorf("couldn't parse giveaway entered points response: %s", err.Error())
	}
	return p, nil
}

type SyncWithSteamResponse struct {
	SyncPrivacyRequirements bool   `json:"sync_privacy_requirements"`
	Type                    string `json:"type"`
	Message                 string `json:"msg"`
}

func ajaxRequest(r steamGiftsRequest) ([]byte, error) {
	client := &http.Client{
		Jar:       newCookieJar(),
		Transport: &http.Transport{},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("POST", steamGiftsUrl+"/ajax.php", strings.NewReader(r.asFormData().Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	log.Println("Sending request...")
	resp, err := client.Do(req)
	log.Println("Done.")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func EnterGiveaway(p *GivewayDetailsPage) (*GivewayEnteredResponse, error) {
	bodyBytes, err := ajaxRequest(p)
	if err != nil {
		return nil, err
	}
	var enteredResp GivewayEnteredResponse
	err = json.Unmarshal(bodyBytes, &enteredResp)
	if err != nil {
		return nil, err
	}
	if enteredResp.RespType == "" {
		return nil, errNotLoggedIn
	}
	return &enteredResp, nil
}

func SyncWithSteam(p *ProfilePage) (*SyncWithSteamResponse, error) {
	bodyBytes, err := ajaxRequest(p)
	if err != nil {
		return nil, err
	}
	var syncedResp SyncWithSteamResponse
	err = json.Unmarshal(bodyBytes, &syncedResp)
	if err != nil {
		return nil, err
	}
	return &syncedResp, nil
}
