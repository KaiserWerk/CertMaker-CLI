package main

import (
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

	// print version and exit
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(AppVersion)
		return
	}

	// determine the cache directory to store the authentication information
	c, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("could not determine app cache directory:", err.Error())
		os.Exit(-1)
	}

	cacheDir := filepath.Join(c, "certctl")
	_ = os.MkdirAll(cacheDir, 0644)

	authFile := filepath.Join(cacheDir, ".auth")

	isInteractive := true
	if len(os.Args) > 1 {
		isInteractive = false
	}

	authenticator := newAuthenticator(authFile)

	if isInteractive {
		err = startInteractiveMode(authenticator)
		if err != nil {
			fmt.Println("interactive mode error:", err.Error())
			os.Exit(-1)

		}
	} else {
		err = startNonInteractiveMode(authenticator, os.Args)
		if err != nil {
			fmt.Println("non-interactive mode error:", err.Error())
			os.Exit(-1)
		}
	}

	//
	//// TODO: set up flags!
	//authEntries, err = setupAuthConfig()
	//if err != nil {
	//	fmt.Println()
	//}

	//pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Sprint("CertMaker CLI"))
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
