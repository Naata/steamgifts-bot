package main

import (
	"log"
	"math/rand/v2"
	"os"
	"sg_bot/internal/config"
	"sg_bot/internal/steamgifts"
	"time"
)

func wait(minMax config.MinMaxSeconds) (int, func()) {
	seconds := rand.IntN(minMax.Max-minMax.Min) + minMax.Min
	return seconds, func() { time.Sleep(time.Duration(seconds) * time.Second) }
}

func enterGiveaways(page *steamgifts.GiveawayListPage, conf *config.Config) {
	log.Printf("Entering %d giveaways: %s", len(page.Hrefs), page.Hrefs)
	for _, href := range page.Hrefs {
		seconds, sleep := wait(conf.WaitForGiveaway)
		log.Printf("Waiting %d seconds before entering giveaway...", seconds)
		sleep()

		details, errors := steamgifts.GetGiveawayDetails(href)
		if errors != nil {
			logErrors(errors)
			continue
		}
		log.Printf("Current user points: %d\n", details.SteamGiftsPage.UserPoints)
		log.Println(details)

		if details.NotEnoughPoints {
			log.Printf("Not enough points to join giveaway '%s', skipping...", details.Name)
			continue
		}
		entered, err := steamgifts.EnterGiveaway(&details)
		if err != nil {
			log.Printf("Couldn't enter giveaway: %s", err.Error())
			continue
		}
		log.Println(entered)
	}
}

func logErrors(errors []error) {
	for _, err := range errors {
		log.Printf("Couldn't parse giveaway details page: %s", err.Error())
	}
}

func syncWithSteam(conf *config.Config) {
	if !conf.SyncWithSteam {
		log.Println("Skipping Steam sync")
		return
	}

	log.Println("Synchronizing with Steam...")
	page, errors := steamgifts.GetProfilePage()
	if errors != nil {
		logErrors(errors)
		return
	}
	r, err := steamgifts.SyncWithSteam(page)
	if err != nil {
		log.Println(err)
	}
	log.Println(r.Message)
}

func main() {
	conf, exists := config.GetConfig()
	if !exists {
		log.Println("User config not fund, default config created.")
		os.Exit(1)
	}

	log.Printf("User config: %s", conf)

	steamgifts.SetSessionId(conf.PhpSessId)
	if !steamgifts.PhpSessIdValid() {
		log.Fatalf("PHPSESSID '%s' is not valid", conf.PhpSessId)
	}

	for {
		syncWithSteam(conf)

		games, errors := steamgifts.GetWishlistedGames()
		if errors != nil {
			log.Println("Error listing wishlist page...")
			logErrors(errors)
			continue
		}
		if len(games.Hrefs) == 0 {
			log.Println("All wishlisted giveaways entered, nothing to do...")
		} else {
			enterGiveaways(&games, conf)
		}

		dlcs, errors := steamgifts.GetDLCs()
		if errors != nil {
			log.Println("Error listing dlcs page...")
			logErrors(errors)
			continue
		}
		if len(dlcs.Hrefs) == 0 {
			log.Println("All dlcs entered, nothing to do...")
		} else {
			enterGiveaways(&dlcs, conf)
		}

		seconds, sleep := wait(conf.WaitForWishlist)
		log.Printf("Waiting %d seconds before accessing wishlist...", seconds)
		sleep()
	}
}
