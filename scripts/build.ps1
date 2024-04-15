$root = $(Resolve-Path "../")
$scripts = $pwd
$dist = "$root/dist"
$src = "$root/cmd/sg_bot"

mkdir $dist -Force
Remove-Item $dist -Recurse -Force

Set-Location $src

$Env:GOOS = "linux"; $Env:GOARM = "5";$Env:GOARCH="arm";go build -o $dist/sg_bot_rpi

$Env:GOOS = "windows";$Env:GOARCH="amd64";go build -o $dist/sg_bot_win.exe

Set-Location $scripts