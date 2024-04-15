package steamgifts

import (
	"encoding/json"
)

func toJsonString(i interface{}) string {
	data, _ := json.MarshalIndent(i, "", "  ")
	return string(data)
}

type parsingErrors struct {
	errors []error
}

type SteamGiftsPage struct {
	UserPoints int
	parsingErrors
}

func (sgp SteamGiftsPage) String() string {
	return toJsonString(sgp)
}

type GivewayDetailsPage struct {
	SteamGiftsPage
	Name            string
	Points          int
	Code            string
	XsrfToken       string
	NotEnoughPoints bool
	parsingErrors
}

func (gd GivewayDetailsPage) String() string {
	return toJsonString(gd)
}

type GiveawayListPage struct {
	Hrefs []string
	SteamGiftsPage
}
