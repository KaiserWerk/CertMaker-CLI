$sourcecode = "main.go"
$target = "build\certctl"
$version = "v0.0.1-alpha"
# Windows, 64-bit
$env:GOOS = 'windows'; $env:GOARCH = 'amd64';             go build -o "$($target)-win64.exe" -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode
# Linux, 64-bit
$env:GOOS = 'linux';   $env:GOARCH = 'amd64';             go build -o "$($target)-linux64"   -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode
# Raspberry Pi
$env:GOOS = 'linux';   $env:GOARCH = 'arm'; $env:GOARM=5; go build -o "$($target)-raspi32"   -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode
# macOS
$env:GOOS = 'darwin';  $env:GOARCH = 'amd64';             go build -o "$($target)-macos64"   -ldflags "-s -w -X 'main.AppVersion=$($version)'" $sourcecode