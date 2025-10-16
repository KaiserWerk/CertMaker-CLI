/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// orderCmd represents the order command
var orderCmd = &cobra.Command{
	Use:   "order",
	Short: "Orders a certificate",
	Long: `Order a new certificate (possibly accompanied by a private key) from the CertMaker instance, 
either using a SimpleRequest or a CSR.`,
	Example: `SimpleRequest: cm order sr --domains example.com,myhost.local --ips 127.0.0.1,192.168.178.1 --days 30
CSR: cm order csr --csrfile /path/to/csr.pem --days 90`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("order called")
	},
}

func init() {
	rootCmd.AddCommand(orderCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// orderCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// orderCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
