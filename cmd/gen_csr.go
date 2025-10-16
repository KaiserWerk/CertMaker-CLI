/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// csrfileCmd represents the csrfile command
var genCSRCmd = &cobra.Command{
	Use:     "csr",
	Short:   "Generates a CSR file",
	Long:    `Generates a Certificate Signing Request (CSR) file based on provided parameters.`,
	Example: `cm gen csr --domains localhost,app.host --ips 158.0.4.5 --keyfile /path/to/key.pem`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gen csr called")
	},
}

func init() {
	genCmd.AddCommand(genCSRCmd)

	genCSRCmd.Flags().StringSliceVar(&domains, "domains", nil, "Comma-separated list of domains for the CSR")
	genCSRCmd.Flags().StringSliceVar(&ips, "ips", nil, "Comma-separated list of IP addresses for the CSR")
	genCSRCmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to the private key file used for signing the CSR")
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
