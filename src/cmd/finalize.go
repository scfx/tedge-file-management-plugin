/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// finalizeCmd represents the finalize command
var finalizeCmd = &cobra.Command{
	Use:   "finalize",
	Short: "Empty function thats called by the sm-agent after remove/install commands",
	Long:  `The finalize command closes a sequence of install and remove commands started by a prepare command.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Printf("Finalize function called but not needed")
	},
}

func init() {
	rootCmd.AddCommand(finalizeCmd)
}
