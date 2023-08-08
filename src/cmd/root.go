/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"file-management/files"
	"log"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
)

var cfgFile string
var config *files.Config
var versions map[string]string
var closeFunc func()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "file-management",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		err := readConfig()
		if err != nil {
			return err
		}
		// Read Versions and Setup Logging
		closeFunc, err = setUpLogs()
		err = readVersions()

		return err
	},

	PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Finished execution of command")
		closeFunc()
		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "/etc/tedge/c8y/file-management.json", "Path to configuration file")

}

func readConfig() error {
	var err error
	config, err = files.ReadConfig(cfgFile)
	return err
}

func readVersions() error {
	var err error
	versions, err = files.ReadVersionFile(config.VersionFile)
	return err
}

func setUpLogs() (func(), error) {
	pathToLogFile := path.Join(config.LogFile, "file-mgmnt-"+time.Now().Format("2006-01-02T15:04:05")+".log")
	logFile, err := os.OpenFile(pathToLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		//Check if path exists and create folder
		if _, err := os.Stat(config.LogFile); os.IsNotExist(err) {
			err := os.MkdirAll(config.LogFile, 0755)
			cobra.CheckErr(err)
		}
		logFile, err = os.OpenFile(pathToLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	log.SetOutput(logFile)
	log.Println("Logging to file: ", pathToLogFile)

	//Setting closing functions
	fn := func() {
		log.Println("Closing log file and reset log-output to stdout")
		logFile.Close()
		log.SetOutput(os.Stdout)
	}
	return fn, nil
}
