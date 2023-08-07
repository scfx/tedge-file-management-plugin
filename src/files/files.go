package files

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Path        string `json:"path"`
	VersionFile string `json:"versionFile"`
	LogFile     string `json:"logFile"`
}

type FileVersion struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ReadConfig(path string) (*Config, error) {
	// Read the TOML file
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	// Unmarshal the JSON data into a Config struct
	var config Config
	if err := json.Unmarshal(byteValue, &config); err != nil {
		return nil, err
	}
	config.VersionFile = filepath.Join(config.Path, config.VersionFile)
	return &config, nil
}

func ReadVersionFile(path string) (map[string]string, error) {
	log.Printf("Reading version file %s", path)
	// Read the JSON file
	jsonData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		//If the file does not exist, return an empty map
		if os.IsNotExist(err) {
			return make(map[string]string), nil
		}
		return nil, err
	}

	// Unmarshal the JSON data into a slice of FileVersion structs
	var fileVersions []FileVersion
	if err := json.Unmarshal(jsonData, &fileVersions); err != nil {
		return nil, err
	}

	// Create a map to access the version of each file
	versions := make(map[string]string)
	for _, fileVersion := range fileVersions {
		versions[fileVersion.Name] = fileVersion.Version
	}
	return versions, nil
}

func WriteVersionFile(path string, versions map[string]string) error {
	// Create a slice of FileVersion structs from the map
	var fileVersions []FileVersion
	for name, version := range versions {
		fileVersions = append(fileVersions, FileVersion{name, version})
	}

	// Marshal the slice of FileVersion structs into JSON
	jsonData, err := json.MarshalIndent(fileVersions, "", "    ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file
	if err := ioutil.WriteFile(path, jsonData, 0644); err != nil {
		return err
	}
	return nil
}

// getFiles returns a list of files in the given path and their versions
func GetFiles(config *Config, versions map[string]string) (map[string]string, error) {
	// Open the directory
	dir, err := os.Open(config.Path)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory contents
	files, err := dir.Readdir(0)
	if err != nil {
		return nil, err
	}

	// Add each file to the map. If the file is not found in the map, set its version to "0.0"
	for _, file := range files {
		//Exclude the conifg.VersionFile file

		if file.Name() != filepath.Base(config.VersionFile) && !file.IsDir() {
			_, found := versions[file.Name()]
			if !found {
				versions[file.Name()] = "0.0"
			}
		}
	}

	return versions, nil
}

func CopyFile(file string, folder string, name string) error {
	// Copy the file to the folder
	// Read the content of the source file
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	// Create the destination file path
	destinationPath := filepath.Join(folder, name)

	// Write the content to the destination file
	err = ioutil.WriteFile(destinationPath, content, 0644)
	if err != nil {
		return err
	}

	return nil

}
