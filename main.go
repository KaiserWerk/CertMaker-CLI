package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	appName = "certctl"
)

var (
	AppVersion = "v0.0.0"
	err        error
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("no arguments provided")
		os.Exit(-1)
	}

	// print version and exit
	if os.Args[1] == "version" {
		fmt.Println(AppVersion)
		return
	}

	// determine the cache directory to store the authentication information
	c, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("could not determine app cache directory:", err.Error())
		os.Exit(-1)
	}

	cacheDir = filepath.Join(c, "certctl")
	_ = os.MkdirAll(cacheDir, 0644)

	authFile := filepath.Join(cacheDir, ".auth")

	isInteractive := false
	if len(os.Args) == 2 && os.Args[1] == "interactive" {
		isInteractive = true
	}

	authenticator := newAuthenticator(authFile)
	err = authenticator.get()
	if err != nil {
		if errors.Is(err, ErrNoAuthFound) {
			fmt.Println("No authentication entries found; please authenticate.")
			return
		} else {
			fmt.Println(err.Error()) // debug
			fmt.Println("An error occurred while fetching authentication info; please authenticate.")
			os.Exit(-1)
		}
	}

	if isInteractive {
		err = startInteractiveMode(authenticator)
	} else {
		err = startNonInteractiveMode(authenticator)
	}

	fmt.Println("is interactive session; not implemented yet")
	//
	//// TODO: set up flags!
	//authEntries, err = setupAuthConfig()
	//if err != nil {
	//	fmt.Println()
	//}

	//pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Sprint("CertMaker Command Line Interface"))
	//pterm.Info.Println("This is the CertMaker Command Line Interface or " + pterm.LightMagenta(appName) + " in short!" +
	//	"\nThe CLI allows you to easily obtain certificates from a running CertMaker" +
	//	"\ninstance and revoking is just as easy!" +
	//	"\n\n" +
	//	"\nIn order to authenticate you will just need a valid API Token." +
	//	"\nYou can actually authenticate against multiple CertMaker instances" +
	//	"\nand choose which one to use at start by choosing the appropriate index." +
	//	"\n\nIf you are authenticated against only one instance, it will" +
	//	"\nautomatically selected at start.")

	//reader := bufio.NewReader(os.Stdin)
	//for {
	//	fmt.Print("> ")
	//	input, err := reader.ReadString('\n')
	//	if err != nil {
	//		fmt.Println("failed to read input:", err.Error())
	//		os.Exit(-1)
	//	}
	//	input = strings.TrimSpace(input)
	//
	//	switch true {
	//	case input == "exit":
	//		fallthrough
	//	case input == "bye":
	//		fmt.Println("Exiting...")
	//		return
	//	case loginN.MatchString(input):
	//		parts := strings.Split(input, " ")
	//		index, _ := strconv.Atoi(parts[1])
	//	case input == "login":
	//		fmt.Print("CertMaker Instance Base URL: ")
	//		var url string
	//		_, _ = fmt.Scanln(&url)
	//		fmt.Print("API Token: ")
	//		var token string
	//		_, _ = fmt.Scanln(&token)
	//
	//		err := authenticate(url, token)
	//	case input == "list instances":
	//	}
	//}
}
