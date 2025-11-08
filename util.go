package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	envXdgDataHome = "XDG_DATA_HOME"
	envXdgDataDirs = "XDG_DATA_DIRS"
	envHome        = "HOME"
)

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

	cursorList := getAssets(cursorSearchPaths, func(fullPath, name string) bool {
		return isFile(filepath.Join(fullPath, "index.theme")) && isDir(filepath.Join(fullPath, "cursors"))
	})

	return cursorList
}

func getAssetSearchPaths(subDir, legacyDir string) []string {
	assetPaths := []string{}

	dataDirs := getDataDirs()

	for _, dir := range dataDirs {
		if isDir(filepath.Join(dir, subDir)) {
			assetPaths = append(assetPaths, filepath.Join(dir, subDir))
		}
	}

	if isDir(filepath.Join(os.Getenv(envHome), legacyDir)) {
		assetPaths = append(assetPaths, filepath.Join(os.Getenv(envHome), legacyDir))
	}

	return assetPaths
}

func getAssets(searchPaths []string, isValidAsset func(fullpath, entryName string) bool) []string {
	assetList := []string{}

	for _, dir := range searchPaths {
		entries, err := os.ReadDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: could not read contents of %s: %s\n", dir, err)
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			fullPath := filepath.Join(dir, entry.Name())

			if isValidAsset(fullPath, entry.Name()) {
				assetList = append(assetList, entry.Name())
			}
		}
	}

	return assetList
}

func getDataDirs() []string {
	systemDataDirs := os.Getenv(envXdgDataDirs)
	if systemDataDirs == "" {
		systemDataDirs = "/usr/local/share:/usr/share"
	}

	homeDataDir := os.Getenv(envXdgDataHome)
	if homeDataDir == "" {
		homeDataDir = filepath.Join(os.Getenv(envHome), ".local", "share")
	}

	dataDirs := strings.Split(systemDataDirs, ":")
	dataDirs = append(dataDirs, homeDataDir)

	return dataDirs
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	if info.IsDir() {
		return false
	}

	return true
}

func printMainHelp(w io.Writer) {
	fmt.Fprintf(w, "Usage: lookctl <command> [arguments]\n\n")
	fmt.Fprintf(w, "Commands:\n")
	fmt.Fprintf(w, "   list		Show installed themes\n")
	fmt.Fprintf(w, "   set		Set the theme, icon, or cursor\n")
	fmt.Fprintf(w, "   current	Show the currently used theme, icon, and cursor\n\n")
	fmt.Fprintf(w, "Run 'lookctl help <command>' for more information on a command.\n")
}

func printListHelp(w io.Writer) {
	fmt.Fprintf(w, "Usage: lookctl list [argument]\n\n")
	fmt.Fprintf(w, "Arguments:\n")
	fmt.Fprintf(w, "   theme	Show installed themes (selected by default)\n")
	fmt.Fprintf(w, "   icon		Show installed icon themes\n")
	fmt.Fprintf(w, "   cursor	Show installed cursor themes\n\n")
}
