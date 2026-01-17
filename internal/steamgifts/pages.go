package steamgifts

import (
	"encoding/json"
	"net/url"
)

type giveawayPage struct {
	name string
	url  string
}

func toJsonString(i interface{}) string {
	data, _ := json.MarshalIndent(i, "", "  ")
	return string(data)
}

type parsingErrors struct {
	errors []error
}

type hasXsrfToken interface {
	xsrfToken() string
	setXsrfToken(t string)
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

func (p *GivewayDetailsPage) xsrfToken() string {
	return p.XsrfToken
}

func (p *GivewayDetailsPage) setXsrfToken(t string) {
	p.XsrfToken = t
}

func (gd GivewayDetailsPage) String() string {
	return toJsonString(gd)
}

func (gd *GivewayDetailsPage) asFormData() *url.Values {
	formData := url.Values{}
	formData.Add("do", "entry_insert")
	formData.Add("xsrf_token", gd.xsrfToken())
	formData.Add("code", gd.Code)
	return &formData
}

type GiveawayListPage struct {
	Name  string
	Hrefs []string
	SteamGiftsPage
}

type ProfilePage struct {
	SteamGiftsPage
	XsrfToken string
}

func (p *ProfilePage) xsrfToken() string {
	return p.XsrfToken
}

func (p *ProfilePage) setXsrfToken(t string) {
	p.XsrfToken = t
}

func (p *ProfilePage) asFormData() *url.Values {
	formData := url.Values{}
	formData.Add("do", "sync")
	formData.Add("xsrf_token", p.xsrfToken())
	return &formData
}

type GiveawayListPageProvider func() (*GiveawayListPage, []error)
