/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/KaiserWerk/CertMaker-CLI/client"

	"github.com/spf13/cobra"
)

var (
	domains  []string
	ips      []string
	emails   []string
	certfile string
	keyfile  string
	days     int
)

// srCmd represents the sr command
var srCmd = &cobra.Command{
	Use:   "sr",
	Short: "Orders a certificate using a SimpleRequest",
	Long: `Order a new certificate accompanied by a private key from the CertMaker instance using a SimpleRequest, meaning you can freely supply
the domain names, IP addresses, email addresses and the desired validity period in days.`,
	Example: `cm order sr --domains example.com,myhost.local --domains newapp.com --ips 192.0.2.1,8.4.1.155 --emails user@example.com --ips 88.77.11.22 --certfile /path/to/cert.pem --keyfile /path/to/key.pem --days 30`,
	Run: func(cmd *cobra.Command, args []string) {

		certData, keyData, err := client.RequestCertificateWithSimpleRequest(domains, ips, emails, days, challengeType)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error requesting certificate:", err)
			return
		}

		err = os.WriteFile(certfile, certData, 0644)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error writing certificate to file:", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Certificate successfully ordered and saved to", certfile)

		err = os.WriteFile(keyfile, keyData, 0600)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error writing private key to file:", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Private key successfully saved to", keyfile)
	},
}

func init() {
	orderCmd.AddCommand(srCmd)

	srCmd.Flags().StringSliceVar(&domains, "domains", nil, "Comma-separated list of domain names. Multiple use allowed.")
	srCmd.Flags().StringSliceVar(&ips, "ips", nil, "Comma-separated list of IP addresses. Multiple use allowed.")
	srCmd.Flags().StringSliceVar(&emails, "emails", nil, "Comma-separated list of email addresses. Multiple use allowed.")
	srCmd.Flags().StringVar(&certfile, "certfile", "", "Path to save the issued certificate (PEM format)")
	srCmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to save the issued private key (PEM format)")

	srCmd.MarkFlagsOneRequired("domains", "ips", "emails")
	srCmd.MarkFlagRequired("certfile")
	srCmd.MarkFlagRequired("keyfile")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// srCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// srCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
