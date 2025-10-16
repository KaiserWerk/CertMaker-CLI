/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/KaiserWerk/CertMaker-CLI/auth"
	"github.com/spf13/cobra"
)

// authremoveCmd represents the authremove command
var authremoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes authentication credentials",
	Long:  `The auth remove command allows users to remove authentication credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := auth.Remove(); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Fprintln(cmd.OutOrStdout(), "nothing to remove")
			} else {
				fmt.Fprintln(cmd.OutOrStderr(), "error removing credentials:", err)
			}
			return
		}
		fmt.Fprintln(cmd.OutOrStdout(), "credentials removed")
	},
}

func init() {
	authCmd.AddCommand(authremoveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authremoveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authremoveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
