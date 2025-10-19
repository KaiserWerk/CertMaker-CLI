/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/pem"
	"fmt"
	"os"

	"github.com/KaiserWerk/CertMaker-CLI/key"
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
		var (
			privKey []byte
			err     error
		)
		switch algo {
		case "rsa":
			privKey, err = key.NewRSA(bits)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error generating RSA key: %v\n", err)
				return
			}
		case "ecdsa":
			privKey, err = key.NewECDSA(bits)
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error generating ECDSA key: %v\n", err)
				return
			}
		case "ed25519":
			privKey, err = key.NewEd25519()
			if err != nil {
				fmt.Fprintf(cmd.OutOrStderr(), "Error generating Ed25519 key: %v\n", err)
				return
			}
		default:
			fmt.Fprintf(cmd.OutOrStderr(), "Unsupported algorithm: %s\n", algo)
			return
		}

		block := &pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privKey,
		}
		privKey = pem.EncodeToMemory(block)

		err = os.WriteFile(keyfile, privKey, 0600)
		if err != nil {
			fmt.Fprintf(cmd.OutOrStderr(), "Error writing key to file: %v\n", err)
			return
		}

		fmt.Fprintf(cmd.OutOrStdout(), "Private key generated and saved to %s\n", keyfile)
	},
}

func init() {
	genCmd.AddCommand(genKeyCmd)

	genKeyCmd.Flags().StringVar(&algo, "algo", "rsa", "The algorithm to use for key generation. Valid values are: rsa, ecdsa, ed25519")
	genKeyCmd.Flags().IntVar(&bits, "bits", 2048, "The key size in bits (for RSA: multiples of 1,024; for ECDSA: 224, 256, 384, 521, ignored for Ed25519)")
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
