/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/KaiserWerk/CertMaker-CLI/auth"
	"github.com/KaiserWerk/CertMaker-CLI/entity"
	"github.com/spf13/cobra"
)

var (
	reason string
)

// revokeCmd represents the revoke command
var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revokes a certificate",
	Long: `Revokes a certificate. This command will communicate with the CertMaker instance to perform the revocation. 
You can only revoke certificates that were issued to you.`,
	Run: func(cmd *cobra.Command, args []string) {
		// read the cert file
		certPEM, err := os.ReadFile(certfile)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error reading certificate file: %v\n", err)
			return
		}

		block, _ := pem.Decode(certPEM)
		if block == nil || block.Type != "CERTIFICATE" {
			fmt.Fprintf(cmd.OutOrStderr(), "Error decoding PEM block from certificate file\n")
			return
		}

		// parse the certificate
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error parsing certificate: %v\n", err)
			return
		}

		revocationRequest := entity.RevocationRequest{
			SerialNumber: cert.SerialNumber.Uint64(),
			Reason:       reason,
		}

		var buf bytes.Buffer
		err = json.NewEncoder(&buf).Encode(revocationRequest)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error encoding revocation request: %v\n", err)
			return
		}

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/api/v1/certificate/revoke", auth.InstanceURL()), &buf)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error creating HTTP request: %v\n", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		auth.SetAuthHeader(req)

		client := &http.Client{Timeout: 30 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error sending HTTP request: %v\n", err)
			return
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		case http.StatusGone:
			fmt.Fprintf(cmd.OutOrStderr(), "Certificate is already revoked.\n")
			return
		case http.StatusBadRequest:
			fmt.Fprintf(cmd.OutOrStderr(), "Bad request: please check your request parameters.\n")
			return
		case http.StatusUnauthorized:
			fmt.Fprintf(cmd.OutOrStderr(), "Unauthorized: please check your authentication credentials.\n")
			return
		case http.StatusForbidden:
			fmt.Fprintf(cmd.OutOrStderr(), "Forbidden: you do not have permission to revoke this certificate or it was not issued to you.\n")
			return
		}

		fmt.Fprintln(cmd.OutOrStdout(), "Certificate successfully revoked.")
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)

	revokeCmd.Flags().StringVar(&certfile, "certfile", "", "The path to the certificate file in PEM-Format.")
	revokeCmd.Flags().StringVar(&reason, "reason", "unspecified", "The reason for revocation. Valid values are: unspecified, keyCompromise, CACompromise, affiliationChanged, superseded, cessationOfOperation, certificateHold, removeFromCRL, privilegeWithdrawn, AACompromise.")
	revokeCmd.MarkFlagsOneRequired("certfile")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// revokeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// revokeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
