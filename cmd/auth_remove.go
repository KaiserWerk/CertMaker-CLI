/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// authremoveCmd represents the authremove command
var authremoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes authentication credentials",
	Long:  `The auth remove command allows users to remove authentication credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("auth remove called")
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
