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

// Flag variables
var filePath string
var moduleVersion string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Installs a software module (file). Usage: install NAME [--module-version VERSION] [--file FILE]",
	Long:  `Installs a software module (file). Usage: install NAME [--module-version VERSION] [--file FILE]`,
	Run: func(cmd *cobra.Command, args []string) {

		// Check if name is provided
		if len(args) < 1 {
			cobra.CheckErr(fmt.Errorf("install needs a name for the command"))
		}
		name := args[0]

		version := moduleVersion
		file := filePath
		// Check if the file and moduleVersion are not empty
		// Check if the file exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			cobra.CheckErr(fmt.Errorf("File to install: '%s' does not exist\n", file))
		}
		// Copy the file to the folder
		err := files.CopyFile(file, config.Path, name)
		cobra.CheckErr(err)
		// Update the version file
		versions[name] = version
		err = files.WriteVersionFile(config.VersionFile, versions)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	// Here you will define your flags and configuration settings.
	installCmd.Flags().StringVarP(&filePath, "file", "f", "", "File to install")
	installCmd.Flags().StringVarP(&moduleVersion, "module-version", "v", "", "Version of the module to install")

	//Mark required flags
	installCmd.MarkFlagRequired("file")
	installCmd.MarkFlagRequired("version")
}
