package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	envXdgDataHome = "XDG_DATA_HOME"
	envXdgDataDirs = "XDG_DATA_DIRS"
	envConfigHome  = "XDG_CONFIG_HOME"
	envHome        = "HOME"
)

func saveCurrentTheme(cfg themeConfig) error {
	if err := saveConfigToFile(cfg); err != nil {
		return err
	}

	if err := saveConfigWithGsettings(cfg); err != nil {
		return err
	}

	return nil
}

func saveConfigToFile(cfg themeConfig) error {
	settingsIniContent := fmt.Sprintf(
		`[Settings]
gtk-theme-name=%s
gtk-icon-theme-name=%s
gtk-cursor-theme-name=%s
gtk-application-prefer-dark-theme=%t`,
		cfg.gtkTheme,
		cfg.iconTheme,
		cfg.cursorTheme,
		cfg.preferDark,
	)

	gtkrc2Content := fmt.Sprintf(
		`gtk-theme-name="%s"
gtk-icon-theme-name="%s"
gtk-cursor-theme-name="%s"`,
		cfg.gtkTheme,
		cfg.iconTheme,
		cfg.cursorTheme,
	)

	configHome := getConfigDir()

	gtk2File := filepath.Join(os.Getenv(envHome), ".gtkrc-2.0")
	gtk3Dir := filepath.Join(configHome, "gtk-3.0")
	gtk4Dir := filepath.Join(configHome, "gtk-4.0")

	if err := os.MkdirAll(gtk3Dir, 0o755); err != nil {
		return fmt.Errorf("failed to create gtk-3.0 directory: %w", err)
	}

	if err := os.MkdirAll(gtk4Dir, 0o755); err != nil {
		return fmt.Errorf("failed to create gtk-4.0 directory: %w", err)
	}

	if err := os.WriteFile(filepath.Join(gtk3Dir, "settings.ini"), []byte(settingsIniContent), 0o644); err != nil {
		return fmt.Errorf("failed to write to gtk-3.0/settings.ini: %w", err)
	}

	if err := os.WriteFile(filepath.Join(gtk4Dir, "settings.ini"), []byte(settingsIniContent), 0o644); err != nil {
		return fmt.Errorf("failed to write to gtk-4.0/settings.ini: %w", err)
	}

	if err := os.WriteFile(gtk2File, []byte(gtkrc2Content), 0o644); err != nil {
		return fmt.Errorf("failed to write to .gtkrc-2.0: %w", err)
	}

	return nil
}

func saveConfigWithGsettings(cfg themeConfig) error {
	colorScheme := "prefer-light"

	if cfg.preferDark {
		colorScheme = "prefer-dark"
	}

	if err := setGsettingsValue(gnomeDesktopInterface, "gtk-theme", cfg.gtkTheme); err != nil {
		return fmt.Errorf("failed to set gtk theme: %w", err)
	}

	if err := setGsettingsValue(gnomeDesktopInterface, "icon-theme", cfg.iconTheme); err != nil {
		return fmt.Errorf("failed to set icon theme: %w", err)
	}

	if err := setGsettingsValue(gnomeDesktopInterface, "cursor-theme", cfg.cursorTheme); err != nil {
		return fmt.Errorf("failed to set cursor theme: %w", err)
	}

	if err := setGsettingsValue(gnomeDesktopInterface, "color-scheme", colorScheme); err != nil {
		return fmt.Errorf("failed to set color scheme: %w", err)
	}

	return nil
}

func setGsettingsValue(schema, key, value string) error {
	cmd := exec.Command("gsettings", "set", schema, key, value)

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func getGsettingsValue(schema, key string) (string, error) {
	cmd := exec.Command("gsettings", "get", schema, key)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	outStr := strings.TrimSpace(string(out))

	outStr = strings.Trim(outStr, "'")

	return outStr, nil
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

func getConfigDir() string {
	configHome := os.Getenv(envConfigHome)
	if configHome == "" {
		configHome = filepath.Join(os.Getenv(envHome), ".config")
	}

	return configHome
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

func containsAny(s string, substrs []string) bool {
	for _, sub := range substrs {
		if strings.Contains(s, sub) {
			return true
		}
	}

	return false
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

func printSetHelp(w io.Writer) {
	fmt.Fprintf(w, "Usage: lookctl set [argument] [theme_name]\n\n")
	fmt.Fprintf(w, "Arguments:\n")
	fmt.Fprintf(w, "   theme	Set theme\n")
	fmt.Fprintf(w, "   icon		Set icon theme\n")
	fmt.Fprintf(w, "   cursor	Set cursor theme\n\n")
}
