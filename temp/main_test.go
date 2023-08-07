package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"os"
	"strings"
	"testing"
)

var configFlag string

func Test_List(t *testing.T) {

	// List without version file should return versions 0.0
	t.Run("List without version file", func(t *testing.T) {
		_, err := setup("./test-folder")
		if err != nil {
			t.Errorf("Error setting up test folder: %s", err)
		}

		defer os.RemoveAll("./test-folder")
		// Run go run main.go
		os.Args = []string{"list"}

		capturedOutputPtr, exitCode := RunAndCaptureOutput(os.Args)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}

		//Check if the output is correct
		expectedOutput := map[string]string{
			"file1.txt": "0.0",
			"file2.txt": "0.0",
		}

		capturedOutputFiles := TransformCapturedOutputToFiles(capturedOutputPtr)

		if len(capturedOutputFiles) != len(expectedOutput) {
			t.Errorf("Expected %d files, got %d", len(expectedOutput), len(capturedOutputFiles))
		}

		for file, version := range expectedOutput {
			if capturedOutputFiles[file] != version {
				t.Errorf("Expected version %s for file %s, got %s", version, file, capturedOutputFiles[file])
			}
		}

	})

	// List with version file should return versions from file
	t.Run("List with version file", func(t *testing.T) {
		versions, err := setup("./test-folder")
		if err != nil {
			t.Errorf("Error setting up test folder: %s", err)
		}

		defer os.RemoveAll("./test-folder")
		//Create a version file for testing and add the versions as json to it
		writeVersionFile("./test-folder/files/.c8y", versions)

		// Run go run main.go
		os.Args = []string{"list"}

		capturedOutputPtr, exitCode := RunAndCaptureOutput(os.Args)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}

		//Check if the output is correct
		expectedOutput := map[string]string{
			"file1.txt": "1.0",
			"file2.txt": "2.0",
		}

		capturedOutputFiles := TransformCapturedOutputToFiles(capturedOutputPtr)

		if len(capturedOutputFiles) != len(expectedOutput) {
			t.Errorf("Expected %d files, got %d", len(expectedOutput), len(capturedOutputFiles))
		}

		for file, version := range expectedOutput {
			if capturedOutputFiles[file] != version {
				t.Errorf("Expected version %s for file %s, got %s", version, file, capturedOutputFiles[file])
			}
		}

	})

}

func Test_Install(t *testing.T) {
	//Install new file to folder and check with list command, if file is added
	t.Run("Install new File to Folder", func(t *testing.T) {
		versions, err := setup("./test-folder")

		//Create a version file for testing and add the versions as json to it
		writeVersionFile("./test-folder/files/.c8y", versions)

		if err != nil {
			t.Errorf("Error setting up test folder: %s", err)
		}

		defer os.RemoveAll("./test-folder")
		//Create dummy file to install in ./test-folder
		file, err := os.Create("./test-folder/file3.txt")
		if err != nil {
			t.Errorf("Error creating dummy file: %s", err)
		}
		file.Close()
		os.Args = []string{"install", "file3.txt", "3.0", "./test-folder/file3.txt"}
		capturedOutputPtr, exitCode := RunAndCaptureOutput(os.Args)
		//Check if exit code is 0
		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}
		//Check if new file is added to folder
		if _, err := os.Stat("./test-folder/files/file3.txt"); os.IsNotExist(err) {
			t.Errorf("Expected file3.txt to be in folder, got %s", err)
		}
		os.Args = []string{"list"}

		capturedOutputPtr, exitCode = RunAndCaptureOutput(os.Args)

		if exitCode != 0 {
			t.Errorf("Expected exit code 0, got %d", exitCode)
		}

		//expected output is versions plus the new file
		expectedOutput := map[string]string{
			"file1.txt": "1.0",
			"file2.txt": "2.0",
			"file3.txt": "3.0",
		}

		capturedOutputFiles := TransformCapturedOutputToFiles(capturedOutputPtr)

		if len(capturedOutputFiles) != len(expectedOutput) {
			t.Errorf("Expected %d files, got %d", len(expectedOutput), len(capturedOutputFiles))
		}

		for file, version := range expectedOutput {
			if capturedOutputFiles[file] != version {
				t.Errorf("Expected version %s for file %s, got %s", version, file, capturedOutputFiles[file])
			}
		}

	})
}

func RunAndCaptureOutput(args []string) (*bytes.Buffer, int) {
	// Override os.Exit during testing
	exitCode := -1
	exitFunc = func(code int) {
		exitCode = code
	}
	defer func() {
		exitFunc = os.Exit
		os.Args = []string{}
	}()
	flagSet := flag.NewFlagSet("test", flag.ExitOnError)
	flag.CommandLine = flagSet

	args = append(args, "-config='./test-folder/config.json'")
	os.Args = args

	// Capture os.Stdout and os.Stderr to check output and error messages
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w

	// Call the function being tested here
	main()

	w.Close()

	// Read the captured output
	var capturedOutput bytes.Buffer
	io.Copy(&capturedOutput, r)

	return &capturedOutput, exitCode
}

func TransformCapturedOutputToFiles(capturedOutput *bytes.Buffer) map[string]string {
	// Split the captured output by newlines to get individual lines
	lines := strings.Split(capturedOutput.String(), "\n")

	capturedOutputFiles := make(map[string]string)

	for _, line := range lines {
		// Split each line by tabs to get Name and Version
		parts := strings.Split(line, "\t")
		if len(parts) == 2 {
			capturedOutputFiles[parts[0]] = parts[1]
		}
	}

	return capturedOutputFiles
}

func setup(path string) (map[string]string, error) {
	// Prepare the test environment
	config := Config{
		Path:        "./test-folder/files",
		VersionFile: ".c8y",
	}
	versions := map[string]string{
		"file1.txt": "1.0",
		"file2.txt": "2.0",
	}
	os.Mkdir(path, 0755)
	os.Mkdir(path+"/files", 0755)

	//Create a config file for testing and add the config as json to it
	configFile, err := os.Create("./test-folder/config.json")
	if err != nil {
		return nil, err
	}
	defer configFile.Close()
	// Marshal the slice of FileVersion structs into JSON
	jsonData, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return nil, err
	}

	// Write the JSON data to the file
	if _, err := configFile.Write(jsonData); err != nil {
		return nil, err
	}

	//Create test files in temp folder
	for name, _ := range versions {
		_, err := os.Create("test-folder/files/" + name)
		if err != nil {
			return nil, err
		}
	}
	return versions, nil
}
