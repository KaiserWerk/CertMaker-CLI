package main

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"flag"
	"fmt"
	certmaker "github.com/KaiserWerk/CertMaker-Go-SDK"
	"os"
	"path/filepath"
	"strings"
)

type AuthEntry struct {
	Name     string `json:"name"`
	BaseURL  string `json:"base_url"`
	Token    string `json:"token"`
	Selected bool   `json:"selected"`
}

const (
	appName = "certctl"
)

var (
	ErrNoAuthSelected = errors.New("certctl: no authentication selected from list")
	ErrNoAuthFound    = errors.New("certctl: no entry file found, empty or malformed")
)

var (
	cacheDir string
	//err         error
	authEntries []*AuthEntry
	currentAuth *AuthEntry
	//loginN      = regexp.MustCompile("login [0-9]+")
)

func main() {
	c, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("could not determine app cache directory:", err.Error())
		os.Exit(-1)
	}

	cacheDir = filepath.Join(c, "certctl")
	_ = os.MkdirAll(cacheDir, 0644)

	currentAuth, err = setupAuthEntries(authEntries, currentAuth)
	if err != nil {
		if errors.Is(err, ErrNoAuthSelected) {
			fmt.Println("Please select a default authentication entry first.")
			os.Exit(-1)
		} else if errors.Is(err, ErrNoAuthFound) {
			fmt.Println("Note: no authentication entries found")
		} else {
			fmt.Println("An error occurred while setting up the authentication entries:", err.Error())
		}
	}

	if len(os.Args) >= 2 {
		subCommand := os.Args[1]
		switch subCommand {
		case "authenticate":
			authFlagSet := flag.NewFlagSet("authenticate", flag.ContinueOnError)
			baseUrl := authFlagSet.String("url", "", "the base URL of the CertMaker instance")
			apikey := authFlagSet.String("key", "", "the API key to use for authentication")
			name := authFlagSet.String("name", "", "a URL-safe unique identifier, e.g. dev, local, int, server01, \"tim's server\"...")
			isDefault := authFlagSet.Bool("default", false, "whether to set the new authentication entry as default")
			err = authFlagSet.Parse(os.Args[2:])
			if err == nil {
				if *baseUrl == "" {
					fmt.Println("missing base URL")
					os.Exit(-1)
				}
				if *apikey == "" {
					fmt.Println("missing API key")
					os.Exit(-1)
				}
				if *name == "" {
					fmt.Println("missing entry name (unique identifier)")
					os.Exit(-1)
				}
			} else {
				fmt.Println("authentication flagset error:", err.Error())
				os.Exit(-1)
			}

			for _, k := range authEntries {
				if k.Name == *name {
					fmt.Println("an entry with this identifier already exists!")
					os.Exit(-1)
				}
			}

			if *isDefault {
				for i, _ := range authEntries {
					authEntries[i].Selected = false
				}
			}

			authEntries = append(authEntries, &AuthEntry{
				Name:     *name,
				BaseURL:  *baseUrl,
				Token:    *apikey,
				Selected: *isDefault,
			})

			err = saveAuthEntries(authEntries)
			if err != nil {
				fmt.Printf("could not save auth entries to file: %s\n", err.Error())
				os.Exit(-1)
			}
		case "use-login":

		case "sr":
			srFlagSet := flag.NewFlagSet("sr", flag.ContinueOnError)
			auth := srFlagSet.String("auth", "", "the authentication entry to use (optional)")
			dnsNames := srFlagSet.String("dnsnames", "", "DNS names/domains")
			ips := srFlagSet.String("ips", "", "IP addresses")
			emails := srFlagSet.String("emails", "", "email addresses")
			days := srFlagSet.Int("days", 7, "The validity in days")
			srOutputDir := srFlagSet.String("out", ".", "the output directory")
			err = srFlagSet.Parse(os.Args[2:])
			if err == nil {
				if *dnsNames == "" && *ips == "" && *emails == "" {
					// create a SimpleRequest and process it
					fmt.Println("one of dnsnames, ips or emails must be set")
					os.Exit(-1)
				}
			} else {
				fmt.Println("sr error:", err)
			}

			if *auth != "" {
				for _, v := range authEntries {
					if v.Name == *auth {
						currentAuth = v
					}
				}
			}

			cache, err := certmaker.NewCache()
			if err != nil {
				fmt.Println("NewCache error:", err.Error())
				os.Exit(-1)
			}
			client := certmaker.NewClient(currentAuth.BaseURL, currentAuth.Token, nil)

			sr := &certmaker.SimpleRequest{
				Domains:        nil,
				IPs:            nil,
				EmailAddresses: nil,
				Days:           *days,
			}

			if *dnsNames != "" {
				sr.Domains = make([]string, 0)
				if !strings.Contains(*dnsNames, ",") {
					sr.Domains = append(sr.Domains, *dnsNames)
				} else {
					parts := strings.Split(*dnsNames, ",")
					for i, v := range parts {
						parts[i] = strings.TrimSpace(v)
					}
					sr.Domains = append(sr.Domains, parts...)
				}
			}

			if *ips != "" {
				sr.IPs = make([]string, 0)
				if !strings.Contains(*ips, ",") {
					sr.IPs = append(sr.IPs, *ips)
				} else {
					parts := strings.Split(*ips, ",")
					for i, v := range parts {
						parts[i] = strings.TrimSpace(v)
					}
					sr.IPs = append(sr.IPs, parts...)
				}
			}

			if *emails != "" {
				sr.EmailAddresses = make([]string, 0)
				if !strings.Contains(*emails, ",") {
					sr.EmailAddresses = append(sr.EmailAddresses, *emails)
				} else {
					parts := strings.Split(*emails, ",")
					for i, v := range parts {
						parts[i] = strings.TrimSpace(v)
					}
					sr.EmailAddresses = append(sr.EmailAddresses, parts...)
				}
			}

			err = client.Request(cache, sr)
			if err != nil {
				fmt.Printf("error fetching certificate: %s\n", err.Error())
				os.Exit(-1)
			}

			fmt.Println("certificate successfully obtained!")

			err = os.Rename(cache.GetCertificatePath(), filepath.Join(*srOutputDir, "cert.pem"))
			if err != nil {
				fmt.Printf("could not move certificate file to output directory: %s\n", err.Error())
				os.Exit(-1)
			}
			err = os.Rename(cache.GetPrivateKeyPath(), filepath.Join(*srOutputDir, "key.pem"))
			if err != nil {
				fmt.Printf("could not move private key file to output directory: %s\n", err.Error())
				os.Exit(-1)
			}
			fmt.Println("moved to output directory!")
			os.Exit(0)
		case "csr":
			csrFlagSet := flag.NewFlagSet("csr", flag.ContinueOnError)
			auth := csrFlagSet.String("auth", "", "the authentication entry to use (optional)")
			csrFile := csrFlagSet.String("file", "", "the CSR file")
			csrOutputDir := csrFlagSet.String("out", ".", "the output directory")
			err = csrFlagSet.Parse(os.Args[2:])
			if err == nil {
				if *csrFile == "" {
					fmt.Println("missing CSR file")
					os.Exit(-1)
				}
			} else {
				fmt.Println("csr error:", err)
				os.Exit(-1)
			}

			data, err := os.ReadFile(*csrFile)
			if err != nil {
				fmt.Println("could not read CSR file:", err.Error())
				os.Exit(-1)
			}

			if *auth != "" {
				for _, v := range authEntries {
					if v.Name == *auth {
						currentAuth = v
					}
				}
			}

			cache, err := certmaker.NewCache()
			if err != nil {
				fmt.Println("NewCache error:", err.Error())
				os.Exit(-1)
			}
			client := certmaker.NewClient(currentAuth.BaseURL, currentAuth.Token, nil)

			b, _ := pem.Decode(data)
			csr, _ := x509.ParseCertificateRequest(b.Bytes)
			err = client.RequestWithCSR(cache, csr)
			if err != nil {
				fmt.Println("certificate request with CSR not successful:", err.Error())
				os.Exit(-1)
			}

			fmt.Println("certificate successfully obtained!")

			err = os.Rename(cache.GetCertificatePath(), filepath.Join(*csrOutputDir, "cert.pem"))
			if err != nil {
				fmt.Printf("could not move certificate file to output directory: %s\n", err.Error())
				os.Exit(-1)
			}
			err = os.Rename(cache.GetPrivateKeyPath(), filepath.Join(*csrOutputDir, "key.pem"))
			if err != nil {
				fmt.Printf("could not move private key file to output directory: %s\n", err.Error())
				os.Exit(-1)
			}
			fmt.Println("moved to output directory!")
		default:
			fmt.Printf("unknown command '%s'\n", subCommand)
			os.Exit(-1)
		}
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

func saveAuthEntries(entries []*AuthEntry) error {
	file := filepath.Join(cacheDir, ".auth")

	j, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	return os.WriteFile(file, j, 0644)
}

func setupAuthEntries(entries []*AuthEntry, currentAuth *AuthEntry) (*AuthEntry, error) {
	file := filepath.Join(cacheDir, ".auth")
	cont, err := os.ReadFile(file)
	if err != nil {
		return nil, ErrNoAuthFound
	}

	err = json.Unmarshal(cont, &entries)
	if err != nil {
		entries = make([]*AuthEntry, 0)
		return nil, ErrNoAuthFound
	}

	if len(entries) == 0 {
		return nil, ErrNoAuthFound
	}
	if len(entries) == 1 {
		return entries[0], nil
	}

	for _, v := range entries {
		if v.Selected {
			return v, nil
		}
	}

	return nil, ErrNoAuthSelected
}
