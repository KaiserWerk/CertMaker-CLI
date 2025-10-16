/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	algo string
	bits int
)

// keyfileCmd represents the keyfile command
var genKeyCmd = &cobra.Command{
	Use:   "key",
	Short: "Generates a private key file",
	Long:  `Generates a private key file in PEM format based on the specified algorithm and key size.`,
	Example: "cm gen key --algo ecdsa --bits 256 --keyfile /path/to/key.pem\n" +
		"cm gen key --algo rsa --bits 2048 --keyfile /path/to/key.pem\n" +
		"cm gen key --algo ed25519 --keyfile /path/to/key.pem",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("keyfile called")
	},
}

func init() {
	genCmd.AddCommand(genKeyCmd)

	genKeyCmd.Flags().StringVar(&algo, "algo", "rsa", "The algorithm to use for key generation. Valid values are: rsa, ecdsa, ed25519")
	genKeyCmd.Flags().IntVar(&bits, "bits", 2048, "The key size in bits (for RSA: multiples of 1,024; for ECDSA: 256, 384, 521, ignored for Ed25519)")
	genKeyCmd.Flags().StringVar(&keyfile, "keyfile", "", "Path to save the generated private key file (PEM format)")
	genKeyCmd.MarkFlagRequired("keyfile")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// keyfileCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// keyfileCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
