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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
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

		logFile, err := os.Create(path.Join(config.LogFile, "file-install-"+time.Now().Format("2006-01-02T15:04:05")+".log"))
		if err != nil {
			//Check if path exists and create folder
			if _, err := os.Stat(config.LogFile); os.IsNotExist(err) {
				err := os.MkdirAll(config.LogFile, 0755)
				cobra.CheckErr(err)
			}
			logFile, err = os.Create(path.Join(config.LogFile, "file-install-"+time.Now().Format("2006-01-02T15:04:05")+".log"))

		}
		cobra.CheckErr(err)
		defer logFile.Close()
		log.SetOutput(logFile)

		//Log config
		log.Printf("Config: Path: %s, versionFile: %s, logFile: %s ", config.Path, config.VersionFile, config.LogFile)
		// Read the version file
		versions, err := files.ReadVersionFile(config.VersionFile)
		cobra.CheckErr(err)

		// List all files in the folder and print their versions
		files, err := files.GetFiles(config, versions)
		cobra.CheckErr(err)
		//Print the files and their versions seperated by a tab
		for file, version := range files {
			fmt.Printf("%s\t%s\n", file, version)
		}
		os.Exit(0)

	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
