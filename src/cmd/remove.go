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
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Read the config file
		config, err := files.ReadConfig(cfgFile)
		cobra.CheckErr(err)
		// Read the version file
		versions, err := files.ReadVersionFile(config.VersionFile)
		cobra.CheckErr(err)

		// Remove a file by deleting it from the folder
		// Check if name is provided
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("remove needs a name for the command"))
		}
		name := args[0]

		// Check if the file exists
		if _, err := os.Stat(config.Path + "/" + name); os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("File to remove '%s' does not exist", name))
		}
		// Delete the file from the folder
		err = os.Remove(config.Path + "/" + name)
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

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// removeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// removeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
