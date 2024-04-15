package steamgifts

import (
	"net/http"

	"github.com/gocolly/colly"
)

const steamGiftsUrl = "https://www.steamgifts.com"
const userAgent = "Mozilla/5.0 (Linux; Android 11) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.0 Safari/537.36"

func collector() *colly.Collector {
	c := colly.NewCollector()
	c.SetCookies(steamGiftsUrl, []*http.Cookie{cookie()})
	return c
}

func PhpSessIdValid() bool {
	c := collector()
	loggedIn := true
	c.OnHTML("a[href='/?login']", func(e *colly.HTMLElement) {
		loggedIn = false
	})
	c.Visit(steamGiftsUrl)
	return loggedIn
}

func scrapeGameListPage(sgpage string) (GiveawayListPage, []error) {
	c := collector()
	page := GiveawayListPage{}
	addWishlistHrefHandler(c, &page)
	addUserPointsHandler(c, &page.SteamGiftsPage)
	c.Visit(steamGiftsUrl + sgpage)
	if len(page.errors) != 0 {
		return page, page.errors
	}
	return page, nil
}

func GetDLCs() (GiveawayListPage, []error) {
	return scrapeGameListPage("/giveaways/search?dlc=true")
}

func GetWishlistedGames() (GiveawayListPage, []error) {
	return scrapeGameListPage("/giveaways/search?type=wishlist")
}

func GetGiveawayDetails(giveawayUrl string) (GivewayDetailsPage, []error) {
	fullUrl := steamGiftsUrl + giveawayUrl
	c := collector()
	page := GivewayDetailsPage{}
	addUserPointsHandler(c, &page.SteamGiftsPage)
	addGiveawayNameHandler(c, &page)
	addGiveawayXsrfTokenHandler(c, &page)
	addGiveawayCodeHandler(c, &page)
	addGiveawayPointsHandler(c, &page)
	addGivewayErrorHandler(c, &page)
	c.Visit(fullUrl)
	if page.XsrfToken == "" && !page.NotEnoughPoints {
		errors := make([]error, 1)
		errors[0] = notLoggedIn
		return page, errors
	}
	if len(page.errors) != 0 {
		return page, page.errors
	}
	return page, nil
}
