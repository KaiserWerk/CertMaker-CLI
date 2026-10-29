$version = git describe --tags --abbrev=0
if ($LASTEXITCODE -ne 0) {
    throw "Git-Tag konnte nicht ermittelt werden."
}

$version = $version.Trim()
if (-not $version) {
    throw "Git hat einen leeren Tag zurückgegeben."
}

$sourcecode = "main.go"
$target = "build\cmctl"

# Linux, 64-bit
$env:GOOS = 'linux';   $env:GOARCH = 'amd64';             go build -o "$($target)-linux64"   -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode
# Raspberry Pi
$env:GOOS = 'linux';   $env:GOARCH = 'arm'; $env:GOARM=5; go build -o "$($target)-raspi32"   -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode
# macOS
$env:GOOS = 'darwin';  $env:GOARCH = 'amd64';             go build -o "$($target)-macos64"   -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode
# Windows, 64-bit
$env:GOOS = 'windows'; $env:GOARCH = 'amd64';             go build -o "$($target)-win64.exe" -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode