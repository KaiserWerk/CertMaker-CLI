/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/KaiserWerk/CertMaker-CLI/auth"
)

var (
	instance string
	apiToken string
)

// authaddCmd represents the authadd command
var authaddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds authentication credentials",
	Long:  `The auth add command allows users to add authentication credentials to be used by 'cm'.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.Set(instance, apiToken); err != nil {
			fmt.Fprintln(cmd.OutOrStderr(), "error adding credentials:", err)
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), "credentials added")
	},
}

func init() {
	authCmd.AddCommand(authaddCmd)

	authaddCmd.Flags().StringVar(&instance, "instance", "", "The Certmaker instance URL (base URL)")
	authaddCmd.Flags().StringVar(&apiToken, "token", "", "The token used for authentication")
	authaddCmd.MarkFlagRequired("instance")
	authaddCmd.MarkFlagRequired("token")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authaddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authaddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
