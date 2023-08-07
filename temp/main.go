package main

import (
	"flag"
	"fmt"
	"os"
	"path"
)

var exitFunc = os.Exit // Custom exit function for testing

func main() {
	// Define flags
	moduleVersionPtr := flag.String("module-version", "", "Version of the module")
	filePtr := flag.String("file", "", "Path to the file")

	// Parse flags
	flag.Parse()

	// Get the non-flag arguments
	args := flag.Args()
	// Read the config file

	fmt.Print(args)
	config, err := readConfig("../file-management.json")

	check_failure(err)
	// Read the version file
	versions, err := readVersionFile(path.Join(config.Path, config.VersionFile))

	check_failure(err)

	// Check if a command was provided
	if len(os.Args) < 2 {
		fmt.Println("Error: No command provided")
		exitFunc(1)
	}
	command := os.Args[1]
	switch command {
	case "list":
		// List all files in the folder and print their versions
		files, err := getFiles(config, versions)
		check_failure(err)
		//Print the files and their versions seperated by a tab
		for file, version := range files {
			fmt.Printf("%s\t%s\n", file, version)
		}
		exitFunc(0)
	case "prepare":
		// No preparation needed

	case "finalize":
		// No finalization needed

	case "install":

		// Install a file by copying it to the folder
		// Check if name, version and file were provided
		if len(os.Args) < 3 {
			fmt.Printf("Error: Not enough arguments provided. Got:%d.", len(os.Args))
			exitFunc(1)
		}

		name := os.Args[2]
		version := *moduleVersionPtr
		file := *filePtr
		// Check if the file and moduleVersion are not empty
		if name == "" || version == "" || file == "" {
			fmt.Printf("Error: Not all flags provided. Name is %s, Version is %s, File is %s", name, version, file)
			exitFunc(1)
		}
		// Check if the file exists
		if _, err := os.Stat(file); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist\n", file)
			exitFunc(1)
		}
		// Copy the file to the folder
		err := copyFile(file, config.Path, name)
		check_failure(err)
		// Update the version file
		versions[name] = version
		err = writeVersionFile(config.VersionFile, versions)
		check_failure(err)
		os.Exit(0)

	case "remove":
		// Remove a file by deleting it from the folder
		// Check if name is provided
		if len(os.Args) < 3 {
			fmt.Println("Error: Not enough arguments provided")
			exitFunc(1)
		}
		name := os.Args[2]
		// Check if the file exists
		if _, err := os.Stat(config.Path + "/" + name); os.IsNotExist(err) {
			fmt.Printf("Error: File '%s' does not exist\n", name)
			exitFunc(1)
		}
		// Delete the file from the folder
		err := os.Remove(config.Path + "/" + name)
		check_failure(err)
		// Update the version file
		delete(versions, name)
		err = writeVersionFile(config.VersionFile, versions)
		check_failure(err)
		exitFunc(0)

	case "update-list":
		// Bulk updates not required
		exitFunc(1)

	default:
		fmt.Printf("Error: Unknown command '%s'\n", command)
		exitFunc(1)

	}

}

func check_failure(err error) {
	if err != nil {
		fmt.Println(err)
		exitFunc(2)
	}
}
