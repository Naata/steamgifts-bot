package config

import (
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	filename                  = "config.json"
	PhpSessId          EnvKey = "SGBOT_PHPSESSID"
	WaitForGiveawayMax EnvKey = "SGBOT_WAITFORGIVEAWAYMAX"
	WaitForGiveawayMin EnvKey = "SGBOT_WAITFORGIVEAWAYMIN"
	SyncWithSteam      EnvKey = "SGBOT_SYNCWITHSTEAM"
	PagesToScan        EnvKey = "SGBOT_PAGESTOSCAN"
)

type MinMaxSeconds struct {
	Min int `json:"min_seconds"`
	Max int `json:"max_seconds"`
}

type Config struct {
	PhpSessId        string        `json:"phpsessid"`
	WaitForGiveaway  MinMaxSeconds `json:"wait_for_giveaway"`
	WaitBetweenScans MinMaxSeconds `json:"wait_between_scans"`
	SyncWithSteam    bool          `json:"sync_with_steam_before_listing"`
	PagesToScan      []string      `json:"pages_to_scan"`
}

func (c Config) String() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}

func defaultConfig() Config {
	return Config{
		PhpSessId:        "",
		WaitForGiveaway:  MinMaxSeconds{Min: 5, Max: 20},
		WaitBetweenScans: MinMaxSeconds{Min: 5 * 60, Max: 15 * 60},
		SyncWithSteam:    true,
		PagesToScan:      []string{"dlc", "wishlist", "multiplecopies", "recommended"},
	}
}

func saveDefaultConfig() {
	config := defaultConfig()
	json, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	err = os.WriteFile(filename, json, fs.ModeAppend)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

func GetConfigFromEnv() (*Config, bool) {
	config := defaultConfig()

	if !PhpSessId.IsSet() {
		panic("PHPSESSID must be set!")
	}
	config.PhpSessId = PhpSessId.GetValue()

	if WaitForGiveawayMax.IsSet() {
		val, err := strconv.Atoi(WaitForGiveawayMax.GetValue())
		if err != nil {
			panic(err)
		}
		config.WaitForGiveaway.Max = val
	}

	if WaitForGiveawayMin.IsSet() {
		val, err := strconv.Atoi(WaitForGiveawayMin.GetValue())
		if err != nil {
			panic(err)
		}
		config.WaitForGiveaway.Min = val
	}

	if SyncWithSteam.IsSet() {
		val, err := strconv.ParseBool(SyncWithSteam.GetValue())
		if err != nil {
			panic(err)
		}
		config.SyncWithSteam = val
	}

	if PagesToScan.IsSet() {
		config.PagesToScan = strings.Split(PagesToScan.GetValue(), ",")
	}

	return &config, true
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
