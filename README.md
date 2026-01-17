# steamgifts-bot
Small Golang tool to automate entering giveaways in SteamGifts page.
It accesses SteamGifts wishlist and DLC pages and enters all giveaways in them.

Accessing wishlist and dlc is configurable and randomized between min and max values.

Entering giveaways is configurable and randomized between min and max values.

## Building application
### For current arch+os
```
cd cmd/sg_bot
go build
```
### For raspberry pi
Set environment variables:
```
GOOS="linux"
GOARM="5"
GOARCH="arm
```
and build like before using `go build`
#### Example for Windows
In powershell console
```powershell
$Env:GOOS = "linux"
$Env:GOARM = "5"
$Env:GOARCH="arm"
go build
```

## Using application
First time running executable creates sample config file, such like:
```json
{
  "phpsessid": "your php session id",
  "wait_for_giveaway": {
    "min_seconds": 5,
    "max_seconds": 20
  },
  "wait_for_wishlist": {
    "min_seconds": 600,
    "max_seconds": 1800
  },
  "sync_with_steam_before_listing": true,
  "pages_to_scan": [
    "dlc",
    "wishlist",
    "multiplecopies",
    "recommended"
  ]
}
```
- `phpsessid` - your PHPSESSID obtained from requests to SteamGifts pages
- `wait_for_giveaway.min_second` - min amount of seconds to wait before entering giveaway
- `wait_for_giveaway.max_second` - max amount of seconds to wait before entering giveaway
- `wait_for_wishlist.min_second` - min amount of seconds to wait before listing wishlist and dlcs
- `wait_for_wishlist.max_second` - max amount of seconds to wait before listing wishlist and dlcs
- `sync_with_steam_before_listing` - sync with steam each time before accessing steamgifts wishlist
- `pages_to_scan` - list of pages to scan, defaults to [ "dlc", "wishlist", "multiplecopies", "recommended" ]

Just paste your `PHPSESSID` and run the binary.
### Obtaining PHPSESSID
1. Log to SteamGifts with your Steam account
1. Open browsers developer tools (on Chrome its `CTRL+SHIFT+I`) 
1. Go to `Network` tab
1. Filter only requests from `steamgifts.com` and refresh page
1. Copy `PHPSESSID` from `Headers` (see screenshot)
1. Paste into `config.json`

Example screenshot:
![phpsessid](docs/readme/phpsessid.png)