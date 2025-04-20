package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	certmaker "github.com/KaiserWerk/CertMaker-Go-SDK"
	"github.com/pterm/pterm"
)

func startInteractiveMode(auth *authenticator) error {
	if err := auth.load(); err != nil {
		for !auth.isAuthenticated() {
			auth.info = askForAuth()
			if err := auth.set(auth.info.BaseURL, auth.info.Token); err != nil {
				return err
			}
		}
	}

	pterm.Info.Println("Current auth:")
	pterm.Info.Println("  " + auth.info.BaseURL)
	pterm.Info.Println("  " + auth.info.Token)
	pterm.Println()

	simpleMode := askForConfirmation("Do you want to use simple mode?")
	if simpleMode {
		pterm.Info.Printfln("Certificate ordering via simple request was selected.")

		sr := certmaker.SimpleRequest{
			Domains:        make([]string, 0, 5),
			IPs:            make([]string, 0, 5),
			EmailAddresses: make([]string, 0, 5),
		}

		dnsNames := make([]string, 0, 5)
		ips := make([]string, 0, 5)
		emails := make([]string, 0, 5)
		for len(dnsNames) == 0 && len(ips) == 0 && len(emails) == 0 {
			dnsNames = askForMultipleStrings("Enter the DNS names for the certificate (one per line)")
			ips = askForMultipleStrings("Enter the IP addresses for the certificate (one per line)")
			emails = askForMultipleStrings("Enter the email addresses for the certificate (one per line)")
		}
		sr.Domains = dnsNames
		sr.IPs = ips
		sr.EmailAddresses = emails
		sr.Days = askForInt("Enter the validity in days for the certificate")

		sr.Subject.Organization = askForString("Enter the organization for the certificate", true)
		sr.Subject.Country = askForString("Enter the country for the certificate (2-letter code)", true)
		sr.Subject.Province = askForString("Enter the province for the certificate", true)
		sr.Subject.Locality = askForString("Enter the locality for the certificate", true)
		sr.Subject.StreetAddress = askForString("Enter the street address for the certificate", true)
		sr.Subject.PostalCode = askForString("Enter the postal code for the certificate", true)

		outputDir := askForString("Enter the output directory for the certificate files", false)
		cache := &certmaker.Cache{
			CacheDir:            outputDir,
			PrivateKeyFilename:  "key.pem",
			CertificateFilename: "cert.pem",
		}
		baseURL, token := auth.authInfo()
		client := certmaker.NewClient(baseURL, token, nil)

		err = client.Request(cache, &sr)
		if err != nil {
			return fmt.Errorf("error fetching certificate: %s\n", err.Error())
		}

		err = os.Rename(cache.GetCertificatePath(), filepath.Join(outputDir, "cert.pem"))
		if err != nil {
			return fmt.Errorf("could not move certificate file to output directory: %s\n", err.Error())
		}
		err = os.Rename(cache.GetPrivateKeyPath(), filepath.Join(outputDir, "key.pem"))
		if err != nil {
			return fmt.Errorf("could not move private key file to output directory: %s\n", err.Error())
		}
		fmt.Println("Certificate files received and written to output directory.")
		return nil
	} else {
		pterm.Info.Printfln("Certificate ordering via CSR was selected.")
		panic("implement me")
	}
}

func askForAuth() *authInfo {
	u := askForString("Please enter the base URL of the CertMaker instance:", false)
	k := askForString("Please enter the token for the CertMaker instance:", false)
	return &authInfo{
		BaseURL: u,
		Token:   k,
	}
}

func askForString(title string, allowEmpty bool) string {
	if title != "" {
		pterm.Info.Printfln(title)
	}

	s, _ := pterm.DefaultInteractiveTextInput.Show(title)

	if s == "" && !allowEmpty {
		pterm.Error.Printfln("Input cannot be empty.")
		return askForString("", allowEmpty)
	}

	return s
}

func askForConfirmation(title string) bool {
	choice, _ := pterm.DefaultInteractiveConfirm.Show(title)
	return choice
}

func askForMultipleStrings(title string) []string {
	textInput := pterm.DefaultInteractiveTextInput.WithMultiLine(true)
	result, _ := textInput.Show(title)
	parts := strings.Split(result, "\n")
	for i := 0; i < len(parts); i++ {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func askForInt(title string) int {
	v := askForString(title, false)
	i, err := strconv.Atoi(v)
	if err != nil {
		pterm.Error.Printfln("Input must be a number.")
		return askForInt("")
	}
	return i
}
