/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	csrfile string
)

// csrCmd represents the csr command
var csrCmd = &cobra.Command{
	Use:   "csr",
	Short: "Orders a certificate using a CSR",
	Long: `Order a new certificate from the CertMaker instance using a Certificate Signing Request (CSR), meaning you provide a CSR file
that contains the necessary information about the certificate.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("csr called")
	},
}

func init() {
	orderCmd.AddCommand(csrCmd)

	csrCmd.Flags().StringVar(&csrfile, "csrfile", "", "Path to the CSR file")
	csrCmd.Flags().IntVar(&days, "days", 7, "Number of days the certificate should be valid for (1-182 days, default 7)")
	csrCmd.MarkFlagRequired("csrfile")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// csrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// csrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
