/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// prepareCmd represents the prepare command
var prepareCmd = &cobra.Command{
	Use:   "prepare",
	Short: "Empty function thats called by the sm-agent before remove/install commands",
	Long:  `The prepare command is invoked by the sm-agent before a sequence of install and remove commands`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Prepare function called but not needed")
	},
}

func init() {
	rootCmd.AddCommand(prepareCmd)
}
