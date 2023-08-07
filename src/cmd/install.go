/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"file-management/files"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
)

// Flag variables
var filePath string
var moduleVersion string

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
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

		// Enable Logging.
		//Filename is install.log + current date and time

		logFile, err := os.Create(path.Join(config.LogFile, "file-install-"+time.Now().Format("2006-01-02T15:04:05")+".log"))
		cobra.CheckErr(err)
		defer logFile.Close()
		log.SetOutput(logFile)

		//Log config
		log.Printf("Config: Path: %s, versionFile: %s, logFile: %s ", config.Path, config.VersionFile, config.LogFile)

		// Read the version file
		versions, err := files.ReadVersionFile(config.VersionFile)
		cobra.CheckErr(err)

		// Remove a file by deleting it from the folder
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
		err = files.CopyFile(file, config.Path, name)
		cobra.CheckErr(err)
		// Update the version file
		versions[name] = version
		err = files.WriteVersionFile(config.VersionFile, versions)
		cobra.CheckErr(err)
		os.Exit(0)
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
