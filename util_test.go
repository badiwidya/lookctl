package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/badiwidya/lookctl/test"
)

func TestGetDataDirs(t *testing.T) {
	tests := []struct {
		description string
		home        string
		xdgDataDir  string
		xdgDataHome string
		want        []string
	}{
		{
			description: "returns correct value if env set",
			home:        "/home/test",
			xdgDataDir:  "/usr/data:/usr/local/data",
			xdgDataHome: "/home/test/data",
			want:        []string{"/usr/data", "/usr/local/data", "/home/test/data"},
		},
		{
			description: "returns fallback value if env not set",
			home:        "/home/test",
			xdgDataDir:  "",
			xdgDataHome: "",
			want:        []string{"/usr/share", "/usr/local/share", "/home/test/.local/share"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			t.Setenv(envHome, tt.home)
			t.Setenv(envXdgDataDirs, tt.xdgDataDir)
			t.Setenv(envXdgDataHome, tt.xdgDataHome)

			got := getDataDirs()

			test.AssertStringSlicesEqual(t, got, tt.want)
		})
	}
}

func TestGetAssertSearchPaths(t *testing.T) {
	tempDir := t.TempDir()

	homePath := filepath.Join(tempDir, "home", "test")
	dataHomePath := filepath.Join(tempDir, homePath, ".local", "share")
	dataDirPath1 := filepath.Join(tempDir, "usr", "share")
	dataDirPath2 := filepath.Join(tempDir, "usr", "local", "share")

	t.Setenv(envHome, homePath)
	t.Setenv(envXdgDataHome, dataHomePath)
	t.Setenv(envXdgDataDirs, fmt.Sprintf("%s:%s", dataDirPath1, dataDirPath2))

	assetPath1 := filepath.Join(homePath, ".themes")
	assetPath2 := filepath.Join(dataHomePath, "themes")
	assetPath3 := filepath.Join(dataDirPath1, "themes")
	// assetPath4 := filepath.Join(dataDirPath2, "themes")

	test.CreateEmptyDir(t, assetPath1)
	test.CreateEmptyDir(t, assetPath2)

	// assetPath3 is a file
	test.CreateEmptyDir(t, dataDirPath1)
	test.CreateEmptyFile(t, assetPath3)

	// assetPath4 nonexist
	test.CreateEmptyDir(t, dataDirPath2)

	got := getAssetSearchPaths("themes", ".themes")

	want := []string{assetPath1, assetPath2}

	test.AssertStringSlicesEqual(t, got, want)
}

func TestGetAssets(t *testing.T) {
	// supress warning output
	ogStderr := os.Stderr

	devNull, _ := os.OpenFile(os.DevNull, os.O_WRONLY, 0)

	os.Stderr = devNull

	t.Cleanup(func() {
		os.Stderr = ogStderr
		devNull.Close()
	})

	tempDir := t.TempDir()

	searchPath1 := filepath.Join(tempDir, "path1")
	searchPath2 := filepath.Join(tempDir, "path2")
	nonexistSearchPath := filepath.Join(tempDir, "nonexist_path")

	test.CreateEmptyDir(t, filepath.Join(searchPath1, "Valid-Theme1"))
	test.CreateEmptyDir(t, filepath.Join(searchPath1, "Invalid-Theme1"))
	test.CreateEmptyDir(t, filepath.Join(searchPath2, "Valid-Theme2"))
	test.CreateEmptyFile(t, filepath.Join(searchPath2, "Invalid-Theme2"))

	searchPaths := []string{searchPath1, searchPath2, nonexistSearchPath}

	mockValidator := func(fullpath, entryName string) bool {
		return strings.HasPrefix(entryName, "Valid-")
	}

	got := getAssets(searchPaths, mockValidator)

	want := []string{"Valid-Theme1", "Valid-Theme2"}

	test.AssertStringSlicesEqual(t, got, want)
}
