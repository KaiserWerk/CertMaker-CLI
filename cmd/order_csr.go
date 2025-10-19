/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/KaiserWerk/CertMaker-CLI/client"

	"github.com/spf13/cobra"
)

var (
	csrfile string
)

// csrCmd represents the csr command
var csrCmd = &cobra.Command{
	Use:   "csr",
	Short: "Orders a certificate using a CSR",
	Long:  `Order a new certificate from the CertMaker instance using a Certificate Signing Request (CSR), which contains the necessary information about the desired certificate.`,
	Run: func(cmd *cobra.Command, args []string) {
		cont, err := os.ReadFile(csrfile)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error reading CSR file:", err)
			return
		}

		block, _ := pem.Decode(cont)
		if block == nil || block.Type != "CERTIFICATE REQUEST" {
			fmt.Fprintln(cmd.OutOrStderr(), "failed to decode PEM block containing CSR or none found")
			return
		}
		_, err = x509.ParseCertificateRequest(block.Bytes)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error parsing CSR:", err)
			return
		}

		certData, err := client.RequestCertificateWithCSR(block.Bytes, days, challengeType)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error requesting certificate:", err)
			return
		}

		block = &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certData,
		}
		certData = pem.EncodeToMemory(block)

		err = os.WriteFile(certfile, certData, 0644)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error writing certificate to file:", err)
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Certificate successfully ordered and saved to", certfile)
	},
}

func init() {
	orderCmd.AddCommand(csrCmd)

	csrCmd.Flags().StringVar(&csrfile, "csrfile", "", "Path to the CSR file")
	csrCmd.Flags().StringVar(&certfile, "certfile", "", "Path where the certificate should be stored")
	csrCmd.MarkFlagRequired("csrfile")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
