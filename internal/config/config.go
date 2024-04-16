package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
)

const filename = "config.json"

type MinMaxSeconds struct {
	Min int `json:"min_seconds"`
	Max int `json:"max_seconds"`
}

type Config struct {
	PhpSessId       string        `json:"phpsessid"`
	WaitForGiveaway MinMaxSeconds `json:"wait_for_giveaway"`
	WaitForWishlist MinMaxSeconds `json:"wait_for_wishlist"`
	SyncWithSteam   bool          `json:"sync_with_steam_before_listing"`
}

func (c Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}

func saveDefaultConfig() {
	config := Config{PhpSessId: "put_php_session_id_here", WaitForGiveaway: MinMaxSeconds{Min: 5, Max: 20}, WaitForWishlist: MinMaxSeconds{Min: 10 * 60, Max: 30 * 60}, SyncWithSteam: true}
	json, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = os.WriteFile(filename, json, fs.ModeAppend)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func GetConfig() (*Config, bool) {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		saveDefaultConfig()
		return nil, false
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Couldn't read config file: " + err.Error())
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Couldn't parse config file: " + err.Error())
	}
	return &config, true
}
