package steamgifts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

func addUserPointsHandler(c *colly.Collector, page *SteamGiftsPage) {
	c.OnHTML("span.nav__points", func(e *colly.HTMLElement) {
		points, err := strconv.Atoi(e.Text)
		if err != nil {
			page.errors = append(page.errors, fmt.Errorf("Couldn't parse user points value: %s", err.Error()))
		}
		page.UserPoints = points
	})
}

func addWishlistHrefHandler(c *colly.Collector, page *GiveawayListPage) {
	c.OnHTML("div.giveaway__row-inner-wrap:not(.is-faded) a.giveaway__heading__name", func(e *colly.HTMLElement) {
		page.Hrefs = append(page.Hrefs, e.Attr("href"))
	})
}

func addGiveawayNameHandler(c *colly.Collector, page *GivewayDetailsPage) {
	c.OnHTML(".featured__heading__medium", func(e *colly.HTMLElement) {
		page.Name = e.Text
	})
}

func addGiveawayXsrfTokenHandler(c *colly.Collector, page *GivewayDetailsPage) {
	c.OnHTML("div.sidebar form input[name='xsrf_token']", func(e *colly.HTMLElement) {
		page.XsrfToken = e.Attr("value")
	})
}

func addGiveawayCodeHandler(c *colly.Collector, page *GivewayDetailsPage) {
	c.OnHTML("div.sidebar form input[name='code']", func(e *colly.HTMLElement) {
		page.Code = e.Attr("value")
	})
}

func addGiveawayPointsHandler(c *colly.Collector, page *GivewayDetailsPage) {
	c.OnHTML("span.sidebar__entry__points", func(e *colly.HTMLElement) {
		s := strings.Replace(e.Text, "(", "", 1)
		s = strings.Replace(s, "P)", "", 1)
		points, err := strconv.Atoi(s)
		if err != nil {
			page.errors = append(page.errors, fmt.Errorf("Couldn't parse giveaways points value: %s", err.Error()))
		}
		page.Points = points
	})
}

func addGivewayErrorHandler(c *colly.Collector, page *GivewayDetailsPage) {
	c.OnHTML("div.sidebar__error.is-disabled", func(e *colly.HTMLElement) {
		page.NotEnoughPoints = true
	})
}
