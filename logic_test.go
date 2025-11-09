package main

import (
	"path/filepath"
	"testing"

	"github.com/badiwidya/lookctl/test"
)

func TestGetInstalledThemes(t *testing.T) {
	tests := []struct {
		description       string
		validThemeNames   []string
		invalidThemeNames []string
		want              []string
	}{
		{
			description:     "returns a list of valid installed themes",
			validThemeNames: []string{"MyTheme", "OurTheme"},
			want:            []string{"MyTheme", "OurTheme"},
		},
		{
			description:       "invalid theme ignored properly",
			invalidThemeNames: []string{"MyTheme", "OurTheme"},
			want:              []string{},
		},
		{
			description:     "excluded themes ignored properly",
			validThemeNames: []string{"Default", "Emacs", "OurTheme"},
			want:            []string{"OurTheme"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			themeDirPath := setupAssetDir(t, "themes")

			for _, validTheme := range tt.validThemeNames {
				themePath := filepath.Join(themeDirPath, validTheme)
				test.CreateEmptyDir(t, themePath)
				test.CreateEmptyFile(t, filepath.Join(themePath, "index.theme"))
			}

			for _, invalidTheme := range tt.invalidThemeNames {
				themePath := filepath.Join(themeDirPath, invalidTheme)
				test.CreateEmptyDir(t, themePath)
			}

			got := getInstalledThemes()

			test.AssertStringSlicesEqual(t, got, tt.want)
		})
	}
}

func TestGetInstalledCursorThemes(t *testing.T) {
	tests := []struct {
		description              string
		validCursorNames         []string
		invalidCursorWithIndex   []string
		invalidCursorWithCursors []string
		want                     []string
	}{
		{
			description:      "returns a list of valid installed cursors",
			validCursorNames: []string{"MyCursor", "OurCursor"},
			want:             []string{"MyCursor", "OurCursor"},
		},
		{
			description:              "icon with cursors subdir but without index.theme ignored properly",
			invalidCursorWithCursors: []string{"InvalidMyCursor", "MyCursor"},
			want:                     []string{},
		},
		{
			description:            "icon with index.theme but without cursors subdir ignored properly",
			invalidCursorWithIndex: []string{"InvalidCursor", "MyCursor"},
			want:                   []string{},
		},
		{
			description:      "excluded cursors ignored properly",
			validCursorNames: []string{"default"},
			want:             []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			iconDirPath := setupAssetDir(t, "icons")

			for _, validCursor := range tt.validCursorNames {
				cursorPath := filepath.Join(iconDirPath, validCursor)
				test.CreateEmptyDir(t, filepath.Join(cursorPath, "cursors"))
				test.CreateEmptyFile(t, filepath.Join(cursorPath, "index.theme"))
			}

			for _, invalidCursor := range tt.invalidCursorWithCursors {
				cursorPath := filepath.Join(iconDirPath, invalidCursor)
				test.CreateEmptyDir(t, filepath.Join(cursorPath, "cursors"))
			}

			for _, invalidCursor := range tt.invalidCursorWithIndex {
				cursorPath := filepath.Join(iconDirPath, invalidCursor)
				test.CreateEmptyDir(t, cursorPath)
				test.CreateEmptyFile(t, filepath.Join(cursorPath, "index.theme"))
			}

			got := getInstalledCursorThemes()

			test.AssertStringSlicesEqual(t, got, tt.want)
		})
	}
}

func TestGetInstalledIconThemes(t *testing.T) {
	tests := []struct {
		description                  string
		validIconNames               []string
		invalidIconWithoutIndex      []string
		invalidIconWithoutIconSubdir []string
		invalidIconIconSubdirAsFile  []string
		want                         []string
	}{
		{
			description:    "returns a list of valid installed icons",
			validIconNames: []string{"MyIcon", "OurIcon"},
			want:           []string{"MyIcon", "OurIcon"},
		},
		{
			description:                  "icon without valid subdir ignored properly",
			validIconNames:               []string{"ValidOne"},
			invalidIconWithoutIconSubdir: []string{"InvalidIcon2", "InvalidIcon1"},
			want:                         []string{"ValidOne"},
		},
		{
			description:             "icon without index.theme ignored properly",
			validIconNames:          []string{"ValidOne"},
			invalidIconWithoutIndex: []string{"InvalidIcon2", "InvalidIcon1"},
			want:                    []string{"ValidOne"},
		},
		{
			description:                 "icon with icon subdir as file ignored properly",
			validIconNames:              []string{"ValidOne"},
			invalidIconIconSubdirAsFile: []string{"ValidButNotIcon"},
			want:                        []string{"ValidOne"},
		},
		{
			description:    "excluded icons ignored properly",
			validIconNames: []string{"ValidOne", "hicolor", "locolor", "default", "gnome"},
			want:           []string{"ValidOne"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			iconDirPath := setupAssetDir(t, "icons")

			for _, validIcon := range tt.validIconNames {
				iconPath := filepath.Join(iconDirPath, validIcon)
				test.CreateEmptyDir(t, filepath.Join(iconPath, "apps"))
				test.CreateEmptyFile(t, filepath.Join(iconPath, "index.theme"))
			}

			for _, invalidIcon := range tt.invalidIconWithoutIndex {
				iconPath := filepath.Join(iconDirPath, invalidIcon)
				test.CreateEmptyDir(t, filepath.Join(iconPath, "apps"))
			}

			for _, invalidIcon := range tt.invalidIconWithoutIconSubdir {
				iconPath := filepath.Join(iconDirPath, invalidIcon)
				test.CreateEmptyDir(t, iconPath)
				test.CreateEmptyFile(t, filepath.Join(iconPath, "index.theme"))
			}

			for _, invalidIcon := range tt.invalidIconIconSubdirAsFile {
				iconPath := filepath.Join(iconDirPath, invalidIcon)
				test.CreateEmptyDir(t, iconPath)
				test.CreateEmptyFile(t, filepath.Join(iconPath, "apps"))
			}

			got := getInstalledIconThemes()

			test.AssertStringSlicesEqual(t, got, tt.want)
		})
	}
}

func setupAssetDir(t testing.TB, assetName string) string {
	t.Helper()

	tempDir := t.TempDir()
	dataDirPath := filepath.Join(tempDir, "usr", "share")

	t.Setenv(envHome, "")
	t.Setenv(envXdgDataHome, "")
	t.Setenv(envXdgDataDirs, dataDirPath)

	assetPath := filepath.Join(dataDirPath, assetName)

	test.CreateEmptyDir(t, assetPath)

	return assetPath
}
