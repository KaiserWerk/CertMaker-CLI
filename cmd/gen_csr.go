/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"
)

// csrfileCmd represents the csrfile command
var genCSRCmd = &cobra.Command{
	Use:     "csr",
	Short:   "Generates a CSR file",
	Long:    `Generates a Certificate Signing Request (CSR) file based on provided parameters.`,
	Example: `cm gen csr --domains localhost,app.host --ips 158.0.4.5 --keyfile /path/to/key.pem --csrfile /path/to/csr.pem`,
	Run: func(cmd *cobra.Command, args []string) {
		privKey, sigAlgo, err := loadPrivateKeyFromFile(keyfile)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error loading private key: %v\n", err)
			return
		}
		tmpl := x509.CertificateRequest{
			DNSNames:           domains,
			IPAddresses:        parseIPs(ips),
			SignatureAlgorithm: sigAlgo,
			// TODO: Subject
		}

		csr, err := x509.CreateCertificateRequest(rand.Reader, &tmpl, privKey)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error creating CSR: %v\n", err)
			return
		}

		csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr})
		err = os.WriteFile(csrfile, csrPEM, 0644)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error saving CSR to file: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "CSR generated and saved to %s\n", csrfile)
	},
}

func init() {
	genCmd.AddCommand(genCSRCmd)

	genCSRCmd.Flags().StringSliceVar(&domains, "domains", nil, "Comma-separated list of domains for the CSR")
	genCSRCmd.Flags().StringSliceVar(&ips, "ips", nil, "Comma-separated list of IP addresses for the CSR")
	genCSRCmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the private key file used for signing the CSR")
	genCSRCmd.Flags().StringVar(&csrfile, "csrfile", "", "Path to save the generated CSR file (PEM format)")
	genCSRCmd.MarkFlagRequired("csrfile")
	genCSRCmd.MarkFlagRequired("keyfile")
	genCSRCmd.MarkFlagsOneRequired("domains", "ips")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csrfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csrfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func parseIPs(ipStrs []string) []net.IP {
	var ips []net.IP
	for _, ipStr := range ipStrs {
		ip := net.ParseIP(ipStr)
		if ip != nil {
			ips = append(ips, ip)
		}
	}
	return ips
}

func loadPrivateKeyFromFile(path string) (crypto.Signer, x509.SignatureAlgorithm, error) {
	pemData, err := os.ReadFile(path)
	if err != nil {
		return nil, x509.UnknownSignatureAlgorithm, err
	}

	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, x509.UnknownSignatureAlgorithm, fmt.Errorf("failed to parse PEM file")
	}

	switch block.Type {
	case "PRIVATE KEY":
		privKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, x509.UnknownSignatureAlgorithm, err
		}
		switch k := privKey.(type) {
		case *rsa.PrivateKey:
			return k, x509.SHA256WithRSA, nil
		case *ecdsa.PrivateKey:
			return k, x509.ECDSAWithSHA256, nil
		case *ed25519.PrivateKey:
			return k, x509.PureEd25519, nil
		}
	case "RSA PRIVATE KEY":
		privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, x509.UnknownSignatureAlgorithm, err
		}
		return privKey, x509.SHA256WithRSA, nil
	}

	return nil, x509.UnknownSignatureAlgorithm, nil
}
