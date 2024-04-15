package steamgifts

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

var phpsessionid string

func cookie() *http.Cookie {
	return &http.Cookie{
		Name:  "PHPSESSID",
		Value: phpsessionid,
	}
}

func newCookieJar() http.CookieJar {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Couldn't create cookie jar: %s", err.Error())
	}
	parsedUrl, err := url.Parse(steamGiftsUrl)
	if err != nil {
		log.Fatalf("Couldn't parse steamgifts url: %s", err.Error())
	}
	jar.SetCookies(parsedUrl, []*http.Cookie{cookie()})
	return jar
}

func SetSessionId(sessionId string) {
	phpsessionid = sessionId
}
