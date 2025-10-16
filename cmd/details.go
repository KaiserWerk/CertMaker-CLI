/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// detailsCmd represents the details command
var detailsCmd = &cobra.Command{
	Use:   "details",
	Short: "Displays detailed information about a certificate",
	Long:  `Displays detailed information about a certificate, including most metadata.`,
	Run: func(cmd *cobra.Command, args []string) {
		cont, err := os.ReadFile(certfile)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error reading certificate file:", err)
			return
		}

		block, _ := pem.Decode(cont)
		if block == nil || block.Type != "CERTIFICATE" {
			fmt.Fprintln(cmd.OutOrStderr(), "failed to decode PEM block containing certificate")
			return
		}
		certificate, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error parsing certificate:", err)
			return
		}

		// Display the certificate details

		fmt.Fprintln(cmd.OutOrStdout(), "Certificate Details:")
		fmt.Fprintln(cmd.OutOrStdout(), "  Subject:", certificate.Subject)
		fmt.Fprintln(cmd.OutOrStdout(), "  Issuer:", certificate.Issuer)
		fmt.Fprintln(cmd.OutOrStdout(), "  Not Before:", certificate.NotBefore)
		fmt.Fprintln(cmd.OutOrStdout(), "  Not After:", certificate.NotAfter)
		fmt.Fprintln(cmd.OutOrStdout(), "  Serial Number:", certificate.SerialNumber)
		fmt.Fprintln(cmd.OutOrStdout(), "  Signature Algorithm:", certificate.SignatureAlgorithm)
		fmt.Fprintln(cmd.OutOrStdout(), "  Public Key Algorithm:", certificate.PublicKeyAlgorithm)
		fmt.Fprintln(cmd.OutOrStdout(), "  Signature (hex):", fmt.Sprintf("%x", certificate.Signature))
		fmt.Fprintln(cmd.OutOrStdout(), "  Version:", certificate.Version)
		fmt.Fprintln(cmd.OutOrStdout(), "  Basic Constraints valid:", certificate.BasicConstraintsValid)
		fmt.Fprintln(cmd.OutOrStdout(), "  Is CA:", certificate.IsCA)
		fmt.Fprintln(cmd.OutOrStdout(), "  Max Path Length:", certificate.MaxPathLen)
		fmt.Fprintln(cmd.OutOrStdout(), "  Max Path Length Zero:", certificate.MaxPathLenZero)

		fmt.Fprintln(cmd.OutOrStdout(), "  Key Usage:")
		if certificate.KeyUsage&x509.KeyUsageDigitalSignature != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Digital Signature")
		}
		if certificate.KeyUsage&x509.KeyUsageContentCommitment != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Content Commitment")
		}
		if certificate.KeyUsage&x509.KeyUsageKeyEncipherment != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Key Encipherment")
		}
		if certificate.KeyUsage&x509.KeyUsageDataEncipherment != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Data Encipherment")
		}
		if certificate.KeyUsage&x509.KeyUsageKeyAgreement != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Key Agreement")
		}
		if certificate.KeyUsage&x509.KeyUsageCertSign != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Certificate Signing")
		}
		if certificate.KeyUsage&x509.KeyUsageCRLSign != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    CRL Signing")
		}
		if certificate.KeyUsage&x509.KeyUsageEncipherOnly != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Encipher Only")
		}
		if certificate.KeyUsage&x509.KeyUsageDecipherOnly != 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "    Decipher Only")
		}

		// Subject Alternative Names (SANs)
		if len(certificate.DNSNames) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  Subject Alternative Names (SAN):")
			for _, dnsName := range certificate.DNSNames {
				fmt.Fprintln(cmd.OutOrStdout(), "  DNS Name:", dnsName)
			}
		}
		if len(certificate.IPAddresses) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  IP Addresses:")
			for _, ip := range certificate.IPAddresses {
				fmt.Fprintln(cmd.OutOrStdout(), "  IP Address:", ip.String())
			}
		}
		if len(certificate.EmailAddresses) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  Email Addresses:")
			for _, email := range certificate.EmailAddresses {
				fmt.Fprintln(cmd.OutOrStdout(), "  Email Address:", email)
			}
		}

		if len(certificate.URIs) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  URIs:")
			for _, uri := range certificate.URIs {
				fmt.Fprintln(cmd.OutOrStdout(), "  URI:", uri.String())
			}
		}

		if len(certificate.PolicyIdentifiers) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  Policy Identifiers:")
			for _, policy := range certificate.PolicyIdentifiers {
				fmt.Fprintln(cmd.OutOrStdout(), "  Policy Identifier:", policy.String())
			}
		}

		if len(certificate.CRLDistributionPoints) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  CRL Distribution Points:")
			for _, crl := range certificate.CRLDistributionPoints {
				fmt.Fprintln(cmd.OutOrStdout(), "  CRL Distribution Point:", crl)
			}
		}

		if len(certificate.OCSPServer) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  OCSP Servers:")
			for _, ocsp := range certificate.OCSPServer {
				fmt.Fprintln(cmd.OutOrStdout(), "  OCSP Server:", ocsp)
			}
		}

		if len(certificate.ExtKeyUsage) > 0 {
			fmt.Fprintln(cmd.OutOrStdout(), "  Extended Key Usages:")
			for _, eku := range certificate.ExtKeyUsage {
				fmt.Fprintln(cmd.OutOrStdout(), "  Extended Key Usage:", eku)
			}
		}

		// if len(certificate.UnknownExtKeyUsage) > 0 {
		// 	fmt.Fprintln(cmd.OutOrStdout(), "  Unknown Extended Key Usages:")
		// 	for _, ueku := range certificate.UnknownExtKeyUsage {
		// 		fmt.Fprintln(cmd.OutOrStdout(), "  Unknown Extended Key Usage:", ueku.String())
		// 	}
		// }

		// if len(certificate.Extensions) > 0 {
		// 	fmt.Fprintln(cmd.OutOrStdout(), "  Extensions:")
		// 	for _, ext := range certificate.Extensions {
		// 		fmt.Fprintln(cmd.OutOrStdout(), "  Extension ID:", ext.Id.String())
		// 		fmt.Fprintln(cmd.OutOrStdout(), "    Critical:", ext.Critical)
		// 		fmt.Fprintln(cmd.OutOrStdout(), "    Value (raw):", ext.Value)
		// 	}
		// }

		// if len(certificate.ExtraExtensions) > 0 {
		// 	fmt.Fprintln(cmd.OutOrStdout(), "  Extra Extensions:")
		// 	for _, ext := range certificate.ExtraExtensions {
		// 		fmt.Fprintln(cmd.OutOrStdout(), "  Extension ID:", ext.Id.String())
		// 		fmt.Fprintln(cmd.OutOrStdout(), "    Critical:", ext.Critical)
		// 		fmt.Fprintln(cmd.OutOrStdout(), "    Value (raw):", ext.Value)
		// 	}
		// }
	},
}

func init() {
	rootCmd.AddCommand(detailsCmd)

	detailsCmd.Flags().StringVar(&certfile, "certfile", "", "Path to the certificate file (PEM format)")
	detailsCmd.MarkFlagRequired("certfile")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// detailsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// detailsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
