package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	certmaker "github.com/KaiserWerk/CertMaker-Go-SDK"
)

func startNonInteractiveMode(auth *authenticator, args []string) error {
	subCommand := args[1]
	switch subCommand {
	case "auth":
		authFlagSet := flag.NewFlagSet("authenticate", flag.ContinueOnError)
		baseUrl := authFlagSet.String("url", "", "the base URL of the CertMaker instance")
		apikey := authFlagSet.String("key", "", "the API key to use for authentication")
		err = authFlagSet.Parse(os.Args[1:])
		if err == nil {
			if *baseUrl == "" {
				fmt.Println("missing base URL")
				os.Exit(-1)
			}
			if *apikey == "" {
				fmt.Println("missing API key")
				os.Exit(-1)
			}

		} else {
			fmt.Println("authentication flagset error:", err.Error())
			os.Exit(-1)
		}

		if err = auth.set(*baseUrl, *apikey); err != nil {
			fmt.Println("could not set authentication information:", err.Error())
			os.Exit(-1)
		}
	case "sr":
		srFlagSet := flag.NewFlagSet("sr", flag.ContinueOnError)

		dnsNames := srFlagSet.String("dnsnames", "", "DNS names/domains")
		ips := srFlagSet.String("ips", "", "IP addresses")
		emails := srFlagSet.String("emails", "", "email addresses")
		days := srFlagSet.Int("days", 90, "The validity in days")
		srOutputDir := srFlagSet.String("out", ".", "the output directory")
		err = srFlagSet.Parse(os.Args[2:])
		if err != nil {
			return fmt.Errorf("failed to parse flags: ", err.Error())
		} else {
			if *dnsNames == "" && *ips == "" && *emails == "" {
				return fmt.Errorf("one of dnsnames, ips or emails must be set")
			}
		}

		if *days < 1 {
			return fmt.Errorf("Days must be greater than or equal to 0.")
		}
		if *days > 365 {
			return fmt.Errorf("Days must be less than or equal to 365.")
		}

		cache, err := certmaker.NewCache()
		if err != nil {
			fmt.Println("NewCache error:", err.Error())
		}
		baseURL, token := auth.authInfo()
		client := certmaker.NewClient(baseURL, token, nil)

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
			return fmt.Errorf("error fetching certificate: %s\n", err.Error())
		}

		err = os.Rename(cache.GetCertificatePath(), filepath.Join(*srOutputDir, "cert.pem"))
		if err != nil {
			return fmt.Errorf("could not move certificate file to output directory: %s\n", err.Error())
		}
		err = os.Rename(cache.GetPrivateKeyPath(), filepath.Join(*srOutputDir, "key.pem"))
		if err != nil {
			return fmt.Errorf("could not move private key file to output directory: %s\n", err.Error())
		}
		fmt.Println("Certificate files received and written to output directory.")
		return nil
	case "csr":
		csrFlagSet := flag.NewFlagSet("csr", flag.ContinueOnError)
		csrFile := csrFlagSet.String("file", "", "the CSR file")
		csrOutputDir := csrFlagSet.String("out", ".", "the output directory")
		err = csrFlagSet.Parse(os.Args[2:])
		if err == nil {
			if *csrFile == "" {
				return fmt.Errorf("missing CSR file")
			}
		} else {
			return fmt.Errorf("csr error:", err)
		}

		data, err := os.ReadFile(*csrFile)
		if err != nil {
			return fmt.Errorf("could not read CSR file:", err.Error())
		}

		cache, err := certmaker.NewCache()
		if err != nil {
			return fmt.Errorf("NewCache error:", err.Error())
		}
		baseURL, token := auth.authInfo()
		client := certmaker.NewClient(baseURL, token, nil)

		b, _ := pem.Decode(data)
		csr, _ := x509.ParseCertificateRequest(b.Bytes)
		err = client.RequestWithCSR(cache, csr)
		if err != nil {
			return fmt.Errorf("certificate request with CSR not successful:", err.Error())
		}

		err = os.Rename(cache.GetCertificatePath(), filepath.Join(*csrOutputDir, "cert.pem"))
		if err != nil {
			return fmt.Errorf("could not move certificate file to output directory: %s\n", err.Error())
		}
		err = os.Rename(cache.GetPrivateKeyPath(), filepath.Join(*csrOutputDir, "key.pem"))
		if err != nil {
			return fmt.Errorf("could not move private key file to output directory: %s\n", err.Error())
		}
		fmt.Println("Certificate files received and written to output directory.")
	default:
		return fmt.Errorf("unknown command '%s'\n", subCommand)
	}

	return nil
}
