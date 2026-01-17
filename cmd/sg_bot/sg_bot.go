package main

import (
	"log"
	"math/rand/v2"
	"os"
	"sg_bot/internal/config"
	"sg_bot/internal/steamgifts"
	"slices"
	"time"
)

func wait(minMax config.MinMaxSeconds) (int, func()) {
	seconds := rand.IntN(minMax.Max-minMax.Min) + minMax.Min
	return seconds, func() { time.Sleep(time.Duration(seconds) * time.Second) }
}

func enterGiveaways(pageProvider steamgifts.GiveawayListPageProvider, conf *config.Config) {
	page, errors := pageProvider()
	if !slices.Contains(conf.PagesToScan, page.Name) {
		log.Printf("Config doesn't request scan of page %s, skipping...", page.Name)
		return
	}

	if errors != nil {
		log.Printf("Error listing %s page...", page.Name)
		logErrors(errors)
		return
	}

	if len(page.Hrefs) == 0 {
		log.Printf("All %s giveaways entered, nothing to do...", page.Name)
		return
	}

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

		enterGiveaways(steamgifts.GetWishlistedGames(), conf)
		enterGiveaways(steamgifts.GetDLCs(), conf)
		enterGiveaways(steamgifts.GetMultipleCopies(), conf)
		enterGiveaways(steamgifts.GetRecommended(), conf)

		seconds, sleep := wait(conf.WaitForWishlist)
		log.Printf("Waiting %d seconds before accessing steamgifts...", seconds)
		sleep()
	}
}
