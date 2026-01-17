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

func scrapeGameListPage(sgpage *giveawayPage) (*GiveawayListPage, []error) {
	c := collector()
	page := GiveawayListPage{Name: sgpage.name}
	addWishlistHrefHandler(c, &page)
	addUserPointsHandler(c, &page.SteamGiftsPage)
	c.Visit(steamGiftsUrl + sgpage.url)
	if len(page.errors) != 0 {
		return &page, page.errors
	}
	return &page, nil
}

func GetDLCs() GiveawayListPageProvider {
	return func() (*GiveawayListPage, []error) {
		return scrapeGameListPage(&giveawayPage{url: "/giveaways/search?dlc=true", name: "dlc"})
	}
}

func GetWishlistedGames() GiveawayListPageProvider {
	return func() (*GiveawayListPage, []error) {
		return scrapeGameListPage(&giveawayPage{url: "/giveaways/search?type=wishlist", name: "wishlist"})
	}
}

func GetMultipleCopies() GiveawayListPageProvider {
	return func() (*GiveawayListPage, []error) {
		return scrapeGameListPage(&giveawayPage{url: "/giveaways/search?copy_min=2", name: "multiplecopies"})
	}
}

func GetRecommended() GiveawayListPageProvider {
	return func() (*GiveawayListPage, []error) {
		return scrapeGameListPage(&giveawayPage{url: "/giveaways/search?type=recommended", name: "recommended"})
	}
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
		errors[0] = errNotLoggedIn
		return page, errors
	}
	if len(page.errors) != 0 {
		return page, page.errors
	}
	return page, nil
}

func GetProfilePage() (*ProfilePage, []error) {
	fullUrl := steamGiftsUrl + "/account/settings/profile"
	page := ProfilePage{}
	c := collector()
	addUserPointsHandler(c, &page.SteamGiftsPage)
	addProfileXsrfTokenHandler(c, &page)
	c.Visit(fullUrl)
	if len(page.errors) != 0 {
		return nil, page.errors
	}

	return &page, nil
}
