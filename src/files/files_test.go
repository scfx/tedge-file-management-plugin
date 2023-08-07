package files

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"
)

// Test readConfig
func TestReadConfig(t *testing.T) {
	// Create a temporary test file
	testFileContent := `{"Path": "/etc/files", "VersionFile": ".c8y"}`
	tmpfile, err := ioutil.TempFile("", "test-config-*.json")
	if err != nil {
		t.Fatalf("Error creating temporary test file: %s", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	if _, err := tmpfile.Write([]byte(testFileContent)); err != nil {
		tmpfile.Close()
		t.Fatalf("Error writing to temporary test file: %s", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temporary test file: %s", err)
	}
	t.Run("Read Config", func(t *testing.T) {
		config, err := ReadConfig(tmpfile.Name())
		if err != nil {
			t.Errorf("Error reading config file: %s", err)
		}
		if config.Path != "/etc/files" {
			t.Errorf("Expected path to be 'files', got '%s'", config.Path)
		}
		if config.VersionFile != path.Join(config.Path, ".c8y") {
			t.Errorf("Expected versionFile to be 'files.json', got '%s'", config.VersionFile)
		}
	})
}

// Test readVersionFile
func TestReadVersionFile(t *testing.T) {
	// Create a temporary test file
	testFileContent := `[{"name": "file1", "version": "1.0"}, {"name": "file2", "version": "2.0"}]`
	tmpfile, err := ioutil.TempFile("", "test-version-*.json")
	if err != nil {
		t.Fatalf("Error creating temporary test file: %s", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	if _, err := tmpfile.Write([]byte(testFileContent)); err != nil {
		tmpfile.Close()
		t.Fatalf("Error writing to temporary test file: %s", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temporary test file: %s", err)
	}
	t.Run("Read Version File", func(t *testing.T) {
		versions, err := ReadVersionFile(tmpfile.Name())
		if err != nil {
			t.Errorf("Error reading version file: %s", err)
		}
		if versions["file1"] != "1.0" {
			t.Errorf("Expected version of file1 to be '1.0', got '%s'", versions["file1"])
		}
		if versions["file2"] != "2.0" {
			t.Errorf("Expected version of file2 to be '2.0', got '%s'", versions["file2"])
		}
	})

	t.Run("Read Empty Version File", func(t *testing.T) {
		// Test when the file does not exist
		nonExistentFile := "non_existent_file.json"
		versions, err := ReadVersionFile(nonExistentFile)
		if err != nil {
			t.Errorf("Error reading version file: %s", err)
		}
		if len(versions) != 0 {
			t.Errorf("Expected an empty version map, but got versions: %v", versions)
		}

	})
}

// Test writeVersionFile
func TestWriteVersionFile(t *testing.T) {
	// Create a temporary test file
	testFileContent := `[{"name": "file1", "version": "1.0"}, {"name": "file2", "version": "2.0"}]`
	tmpfile, err := ioutil.TempFile("", "test-version-*.json")
	if err != nil {
		t.Fatalf("Error creating temporary test file: %s", err)
	}
	defer os.Remove(tmpfile.Name()) // Clean up the temporary file

	if _, err := tmpfile.Write([]byte(testFileContent)); err != nil {
		tmpfile.Close()
		t.Fatalf("Error writing to temporary test file: %s", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatalf("Error closing temporary test file: %s", err)
	}

	t.Run("Write Version File", func(t *testing.T) {
		versions := map[string]string{
			"file1": "1.0",
			"file2": "2.1",
			"file3": "3.0",
		}
		err := WriteVersionFile(tmpfile.Name(), versions)
		if err != nil {
			t.Errorf("Error writing version file: %s", err)
		}

		// Read the file again and check that the content is correct
		versions, err = ReadVersionFile(tmpfile.Name())
		if err != nil {
			t.Errorf("Error reading version file: %s", err)
		}
		if versions["file1"] != "1.0" {
			t.Errorf("Expected version of existing file1 to be '1.0', got '%s'", versions["file1"])
		}
		if versions["file2"] != "2.1" {
			t.Errorf("Expected version of updated file2 to be '2.1', got '%s'", versions["file2"])
		}
		if versions["file3"] != "3.0" {
			t.Errorf("Expected version of newly added file3 to be '3.0', got '%s'", versions["file3"])
		}
	})
}

// Test getFiles
func TestGetFiles(t *testing.T) {
	// Create a temporary test directory
	tmpdir, err := ioutil.TempDir("", "test-files2-*")
	if err != nil {
		t.Fatalf("Error creating temporary test directory: %s", err)
	}
	defer os.RemoveAll(tmpdir) // Clean up the temporary directory

	// Create 2 temporary test files
	var tempFiles []*os.File
	for i := 0; i < 2; i++ {
		// Create a temporary test files
		tmpfile, err := ioutil.TempFile(tmpdir, "test-file-*.txt")
		if err != nil {
			t.Fatalf("Error creating temporary test file: %s", err)
		}
		tempFiles = append(tempFiles, tmpfile)
	}
	test_versions := make(map[string]string)
	for _, fileVersion := range tempFiles {
		test_versions[filepath.Base(fileVersion.Name())] = "1.0"
	}

	t.Run("Get Files", func(t *testing.T) {
		files, err := GetFiles(&Config{Path: tmpdir, VersionFile: ".c8y"}, test_versions)
		if err != nil {
			t.Errorf("Error getting files: %s", err)
		}
		// Check that the correct number of files were returned
		if len(files) != len(tempFiles) {
			t.Errorf("Expected %d files to be returned, but got %d", len(tempFiles), len(files))
		}

		// Check that the correct files were returned
		for _, file := range tempFiles {
			if _, ok := files[filepath.Base(file.Name())]; !ok {
				t.Errorf("Expected file '%s' to be returned, but it was not", file.Name())
			}
		}

		// Check that returned files have the correct version
		for _, file := range tempFiles {
			if files[filepath.Base(file.Name())] != "1.0" {
				t.Errorf("Expected version of file '%s' to be '1.0', got '%s'", file.Name(), files[filepath.Base(file.Name())])
			}
		}

		// Add a new file without versioning
		tmp_noversion, err := ioutil.TempFile(tmpdir, "test-file-*.txt")
		if err != nil {
			t.Fatalf("Error creating temporary test file: %s", err)
		}
		// Get files again
		files, err = GetFiles(&Config{Path: tmpdir, VersionFile: ".c8y"}, test_versions)
		if err != nil {
			t.Errorf("Error getting files: %s", err)
		}
		// check if new file was returned
		if _, ok := files[filepath.Base(tmp_noversion.Name())]; !ok {
			t.Errorf("Expected file '%s' to be returned, but it was not", tmp_noversion.Name())
		}
		// check if new file has default version
		if files[filepath.Base(tmp_noversion.Name())] != "0.0" {
			t.Errorf("Expected version of file '%s' to be '0.0', got '%s'", tmp_noversion.Name(), files[filepath.Base(tmp_noversion.Name())])
		}

	})

	// Test that version file is not returned
	t.Run("Get Files - Version File should not be included", func(t *testing.T) {
		ioutil.WriteFile(path.Join(tmpdir, ".c8y"), []byte("1.0"), 0644)

		//Get Files
		files, err := GetFiles(&Config{Path: tmpdir, VersionFile: path.Join(tmpdir, ".c8y")}, map[string]string{})
		if err != nil {
			t.Errorf("Error getting files: %s", err)
		}

		// Check that .c8y is not in returned list
		if _, ok := files[".c8y"]; ok {
			t.Errorf("Expected file '.c8y' to not be returned, but it was")
		}
	})

}

// Test copyFiles
func TestCopyFiles(t *testing.T) {
	// Create a temporary destination directory
	tmpdest, err := ioutil.TempDir("", "test-files2-*")
	if err != nil {
		t.Fatalf("Error creating temporary test directory: %s", err)
	}
	defer os.RemoveAll(tmpdest) // Clean up the temporary directory

	// Create a temporary source directory
	tmpsrc, err := ioutil.TempDir("", "test-files2-*")
	if err != nil {
		t.Fatalf("Error creating temporary test directory: %s", err)
	}
	defer os.RemoveAll(tmpsrc) // Clean up the temporary directory

	// Create a temporary test file
	tmpfile, err := ioutil.TempFile(tmpsrc, "test-file-*.txt")
	if err != nil {
		t.Fatalf("Error creating temporary test file: %s", err)
	}

	t.Run("Copy Files", func(t *testing.T) {
		// Copy tmpfile from tempdest to tempsrc
		err := CopyFile(filepath.Join(tmpsrc, filepath.Base(tmpfile.Name())), tmpdest, filepath.Base(tmpfile.Name()))
		if err != nil {
			t.Errorf("Error copying files: %s", err)
		}

		// Check that the file was copied
		if _, err := os.Stat(filepath.Join(tmpdest, filepath.Base(tmpfile.Name()))); os.IsNotExist(err) {
			t.Errorf("Expected file '%s' to be copied to '%s', but it was not", tmpfile.Name(), tmpdest)
		}

		// Check that the file has the correct content

	})
}
