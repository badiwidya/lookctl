package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

const gnomeDesktopInterface = "org.gnome.desktop.interface"

type themeConfig struct {
	gtkTheme    string
	iconTheme   string
	cursorTheme string
}

func getCurrentTheme() (themeConfig, error) {
	gtkTheme, err := getGsettingsValue(gnomeDesktopInterface, "gtk-theme")
	if err != nil {
		return themeConfig{}, fmt.Errorf("error: failed to read gtk theme information: %w", err)
	}

	iconTheme, err := getGsettingsValue(gnomeDesktopInterface, "icon-theme")
	if err != nil {
		return themeConfig{}, fmt.Errorf("error: failed to read icon theme information: %w", err)
	}

	cursorTheme, err := getGsettingsValue(gnomeDesktopInterface, "cursor-theme")
	if err != nil {
		return themeConfig{}, fmt.Errorf("error: failed to read cursor theme information: %w", err)
	}

	return themeConfig{
		gtkTheme:    gtkTheme,
		iconTheme:   iconTheme,
		cursorTheme: cursorTheme,
	}, nil
}

func getInstalledThemes() []string {
	themeSearchPaths := getAssetSearchPaths("themes", ".themes")

	excluded := []string{"Default", "Emacs"}

	themeList := getAssets(themeSearchPaths, func(fullPath, name string) bool {
		if slices.Contains(excluded, name) {
			return false
		}

		return isFile(filepath.Join(fullPath, "index.theme"))
	})

	return themeList
}

func getInstalledIconThemes() []string {
	iconSearchPaths := getAssetSearchPaths("icons", ".icons")

	mustContainOneOf := []string{
		"scalable", "apps", "16x16",
		"22x22", "24x24", "32x32",
		"36x36", "48x48", "64x64",
		"72x72", "96x96", "128x128",
		"256x256", "512x512", "mimetypes",
	}

	excluded := []string{"hicolor", "locolor", "default", "gnome"}

	iconList := getAssets(iconSearchPaths, func(fullPath, name string) bool {
		if slices.Contains(excluded, name) {
			return false
		}

		if !isFile(filepath.Join(fullPath, "index.theme")) {
			return false
		}

		content, err := os.ReadDir(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not read contents of %s: %s\n", fullPath, err)

			return false
		}

		return slices.ContainsFunc(content, func(d os.DirEntry) bool {
			return d.IsDir() && slices.Contains(mustContainOneOf, d.Name())
		})
	})

	return iconList
}

func getInstalledCursorThemes() []string {
	cursorSearchPaths := getAssetSearchPaths("icons", ".icons")

	excluded := []string{"default"}

	cursorList := getAssets(cursorSearchPaths, func(fullPath, name string) bool {
		if slices.Contains(excluded, name) {
			return false
		}

		return isFile(filepath.Join(fullPath, "index.theme")) && isDir(filepath.Join(fullPath, "cursors"))
	})

	return cursorList
}
