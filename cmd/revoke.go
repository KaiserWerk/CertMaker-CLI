/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	sn     string
	reason string
)

// revokeCmd represents the revoke command
var revokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revokes a certificate",
	Long: `Revokes a certificate. This command will communicate with the CertMaker instance to perform the revocation. 
You can only revoke certificates that were issued to you. The --certfile flag takes precedence over the --sn flag if both are provided.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Fprintln(cmd.OutOrStdout(), "revoke: not implemented yet")
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)

	revokeCmd.Flags().StringVar(&sn, "sn", "", "Serial number of the certificate to be revoked (hexadecimal, e.g., '01ab23cd'). If the --certfile flag is supplied as well, it is used instead.")
	revokeCmd.Flags().StringVar(&certfile, "certfile", "", "The path to the certificate file in PEM-Format. If the --sn flag is supplied as well, it is ignored.")
	revokeCmd.Flags().StringVar(&reason, "reason", "unspecified", "The reason for revocation. Valid values are: unspecified, keyCompromise, CACompromise, affiliationChanged, superseded, cessationOfOperation, certificateHold, removeFromCRL, privilegeWithdrawn, AACompromise.")
	revokeCmd.MarkFlagsMutuallyExclusive("sn", "certfile")
	revokeCmd.MarkFlagsOneRequired("sn", "certfile")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// revokeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// revokeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
