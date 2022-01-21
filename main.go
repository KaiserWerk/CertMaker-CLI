package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type AuthEntry struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

var (
	err         error
	authEntries []AuthEntry
	loginN      = regexp.MustCompile("login [0-9]+")
)

func main() {
	// TODO: set up flags!
	authEntries, err = setupAuthConfig()
	if err != nil {
		fmt.Println()
	}

	//pterm.DefaultCenter.Print(pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Sprint("CertMaker Command Line Interface"))
	//pterm.Info.Println("This is the CertMaker Command Line Interface or " + pterm.LightMagenta("CertMaker CLI") + " in short!" +
	//	"\nThe CLI allows you to easily obtain certificates from a running CertMaker" +
	//	"\ninstance and revoking is just as easy!" +
	//	"\n\n" +
	//	"\nIn order to authenticate you will just need a valid API Token." +
	//	"\nYou can actually authenticate against multiple CertMaker instances" +
	//	"\nand choose which one to use at start by choosing the appropriate index." +
	//	"\n\nIf you are authenticated against only one instance, it will" +
	//	"\nautomatically selected at start.")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("failed to read input:", err.Error())
			os.Exit(-1)
		}
		input = strings.TrimSpace(input)

		switch true {
		case input == "exit":
			fallthrough
		case input == "bye":
			fmt.Println("Exiting...")
			return
		case loginN.MatchString(input):
			parts := strings.Split(input, " ")
			index, _ := strconv.Atoi(parts[1])
		case input == "login":
			fmt.Print("CertMaker Instance Base URL: ")
			var url string
			_, _ = fmt.Scanln(&url)
			fmt.Print("API Token: ")
			var token string
			_, _ = fmt.Scanln(&token)

			err := authenticate(url, token)
		case input == "list instances":
		}
	}
}

func setupAuthConfig() ([]AuthEntry, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	authDir := filepath.Join(homeDir, "CertMakerCLI")
	if err := os.Mkdir(authDir, 0644); err != nil {
		return nil, err
	}
	file := filepath.Join(authDir, "entries.json")

	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
		return nil, nil
	}

	cont, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if len(cont) == 0 {
		return nil, nil
	}

	var e []AuthEntry
	err = json.Unmarshal(cont, &e)
	if err != nil {
		return nil, err
	}

	return e, nil
}

func setupHTTPClient(url string, token string) error {

}
