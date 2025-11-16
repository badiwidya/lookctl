package main

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const gnomeDesktopInterface = "org.gnome.desktop.interface"

type themeConfig struct {
	gtkTheme    string
	iconTheme   string
	cursorTheme string
	preferDark  bool
}

func getCurrentTheme() (themeConfig, error) {
	gtkTheme, err := getGsettingsValue(gnomeDesktopInterface, "gtk-theme")
	if err != nil {
		return themeConfig{}, fmt.Errorf("failed to read gtk theme information: %w", err)
	}

	iconTheme, err := getGsettingsValue(gnomeDesktopInterface, "icon-theme")
	if err != nil {
		return themeConfig{}, fmt.Errorf("failed to read icon theme information: %w", err)
	}

	cursorTheme, err := getGsettingsValue(gnomeDesktopInterface, "cursor-theme")
	if err != nil {
		return themeConfig{}, fmt.Errorf("failed to read cursor theme information: %w", err)
	}

	colorScheme, err := getGsettingsValue(gnomeDesktopInterface, "color-scheme")
	if err != nil {
		return themeConfig{}, fmt.Errorf("failed to read color scheme information: %w", err)
	}

	preferDark := false
	if colorScheme == "prefer-dark" {
		preferDark = true
	}

	return themeConfig{
		gtkTheme:    gtkTheme,
		iconTheme:   iconTheme,
		cursorTheme: cursorTheme,
		preferDark:  preferDark,
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

	slices.Sort(themeList)

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

	slices.Sort(iconList)

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

	slices.Sort(cursorList)

	return cursorList
}

func setTheme(themeName string) error {
	installedThemes := getInstalledThemes()

	if !slices.Contains(installedThemes, themeName) {
		return fmt.Errorf("theme not found. see 'lookctl list theme' for list available themes")
	}

	currentConfig, err := getCurrentTheme()
	if err != nil {
		return err
	}

	lightKeywords := []string{
		"light",
		"snow",
		"white",
	}

	darkKeywords := []string{
		"dark",
		"dracula",
		"gruvbox",
		"nord",
		"night",
	}

	lname := strings.ToLower(themeName)

	var preferDark bool
	if containsAny(lname, lightKeywords) {
		preferDark = false
	} else if containsAny(lname, darkKeywords) {
		preferDark = true
	}

	currentConfig.gtkTheme = themeName
	currentConfig.preferDark = preferDark

	if err := saveCurrentTheme(currentConfig); err != nil {
		return err
	}

	return nil
}

func setIconTheme(themeName string) error {
	installedIconThemes := getInstalledIconThemes()

	if !slices.Contains(installedIconThemes, themeName) {
		return fmt.Errorf("icon theme not found. see 'lookctl list icon' for list available themes")
	}

	currentConfig, err := getCurrentTheme()
	if err != nil {
		return err
	}

	currentConfig.iconTheme = themeName

	if err := saveCurrentTheme(currentConfig); err != nil {
		return err
	}

	return nil
}

func setCursorTheme(themeName string) error {
	installedCursorThemes := getInstalledCursorThemes()

	if !slices.Contains(installedCursorThemes, themeName) {
		return fmt.Errorf("cursor theme not found. see 'lookctl list cursor' for list available themes")
	}

	currentConfig, err := getCurrentTheme()
	if err != nil {
		return err
	}

	currentConfig.cursorTheme = themeName

	if err := saveCurrentTheme(currentConfig); err != nil {
		return err
	}

	return nil
}
