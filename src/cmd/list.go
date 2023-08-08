/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"file-management/files"
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all files installed",
	Long:  `Returns the list of software modules that have been installed with this plugin, using tab separated values`,
	Run: func(cmd *cobra.Command, args []string) {

		// List all files in the folder and print their versions
		files, err := files.GetFiles(config, versions)
		cobra.CheckErr(err)
		//Print the files and their versions seperated by a tab
		for file, version := range files {
			fmt.Printf("%s\t%s\n", file, version)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
