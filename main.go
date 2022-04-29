package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	certmaker "github.com/KaiserWerk/CertMaker-Go-SDK"
)

type AuthEntry struct {
	Name    string `json:"name"`
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

const (
	appName = "certctl"
)

var (
	err         error
	authEntries map[string]AuthEntry
	currentAuth *AuthEntry
	loginN      = regexp.MustCompile("login [0-9]+")
)

func main() {
	err = setupAuthEntries(authEntries)
	subCommand := os.Args[1]
	switch subCommand {
	case "authenticate":

	case "use-login":

	case "sr":
		srFlagSet := flag.NewFlagSet("sr", flag.ContinueOnError)
		dnsNames := srFlagSet.String("dnsnames", "", "DNS names/domains")
		ips := srFlagSet.String("ips", "", "IP addresses")
		emails := srFlagSet.String("emails", "", "email addresses")
		srOutputDir := srFlagSet.String("out", ".", "the output directory")
		err = srFlagSet.Parse(os.Args[2:])
		if err == nil {
			if *dnsNames != "" || *ips != "" || *emails != "" {
				// create a SimpleRequest and process it
				fmt.Println("it's a sr")
				fmt.Println("dns names:", *dnsNames)
				fmt.Println("ips:", *ips)
				fmt.Println("emails:", *emails)
				fmt.Println("out:", *srOutputDir)
				os.Exit(-1)
			}
		} else {
			fmt.Println("sr error:", err)
		}

		if currentAuth != nil {
			currentAuth = getCurrentAuth()
		}

		cache, err := certmaker.NewCache()
		if err != nil {
			fmt.Println("NewCache error:", err.Error())
			os.Exit(-1)
		}
		client := certmaker.NewClient(currentAuth.BaseURL, currentAuth.Token, nil)
	case "csr":
		csrFlagSet := flag.NewFlagSet("csr", flag.ContinueOnError)
		csrFile := csrFlagSet.String("file", "", "the CSR file")
		csrOutputDir := csrFlagSet.String("out", ".", "the output directory")
		err = csrFlagSet.Parse(os.Args[1:])
		if err == nil {
			fmt.Println("it's a csr")
			fmt.Println("csr file:", *csrFile)
			fmt.Println("out:", *csrOutputDir)
		} else {
			fmt.Println("csr error:", err)
		}
	default:
		fmt.Printf("unknown command '%s'\n", subCommand)
		os.Exit(-1)
	}

	os.Exit(-10)
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

func setupAuthEntries(entries map[string]AuthEntry) error {
	
}

func getCurrentAuth() *AuthEntry {
	return &AuthEntry{ // for debugging purposes only!
		BaseURL: "http://localhost:8880",
		Token:   "9a2c749d35902ee06dcd5ee4fc364434293da4a0b267297bc9a3bf17d8e68e9fcc03d3683a484565",
	}
}

func requestViaSimpleRequest(sr *certmaker.SimpleRequest) error {
	return nil
}

//func setupAuthConfig() ([]AuthEntry, error) {
//	homeDir, err := os.UserHomeDir()
//	if err != nil {
//		return nil, err
//	}
//	authDir := filepath.Join(homeDir, appName)
//	if err := os.Mkdir(authDir, 0644); err != nil {
//		return nil, err
//	}
//	file := filepath.Join(authDir, "entries.json")
//
//	if _, err := os.Stat(file); err != nil && os.IsNotExist(err) {
//		return nil, nil
//	}
//
//	cont, err := ioutil.ReadFile(file)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(cont) == 0 {
//		return nil, nil
//	}
//
//	var e []AuthEntry
//	err = json.Unmarshal(cont, &e)
//	if err != nil {
//		return nil, err
//	}
//
//	return e, nil
//}
//
//func setupHTTPClient(url string, token string) error {
//
//}
