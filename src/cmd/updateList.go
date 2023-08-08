/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// updateListCmd represents the updateList command
var updateListCmd = &cobra.Command{
	Use:   "update-list",
	Short: "Returns Exit Code 1 so that install/remove will be called individually by the sm-agent",
	Long:  `The update-list command accepts a list of software modules and associated operations as install or remove. Returns Exit Code 1 so that install/remove will be called individually by the sm-agent`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Update List function called. Will return Exit Code 1, so that install/remove will be called individually")
		os.Exit(1)
	},
}

func init() {
	rootCmd.AddCommand(updateListCmd)
}
