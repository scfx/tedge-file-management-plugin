/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"file-management/files"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Uninstalls a software module (file). Usage: remove NAME",
	Long:  `Uninstalls a software module (file). Usage: remove NAME`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("remove needs a name for the command"))
		}
		name := args[0]

		// Check if the file exists
		if _, err := os.Stat(config.Path + "/" + name); os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("File to remove '%s' does not exist", name))
		}
		// Delete the file from the folder
		err := os.Remove(config.Path + "/" + name)
		cobra.CheckErr(err)
		// Update the version file
		delete(versions, name)
		err = files.WriteVersionFile(config.VersionFile, versions)
		cobra.CheckErr(err)
		os.Exit(0)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
