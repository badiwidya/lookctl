package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/tabwriter"
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

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	return fs
}

func parseFlag(fs *flag.FlagSet, args []string, printHelp func(*tabwriter.Writer)) error {
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printHelp(newTabWriter(os.Stdout))
			os.Exit(0)
		}

		if strings.Contains(err.Error(), "provided but not defined") {
			flagName := strings.Split(err.Error(), "-")
			return fmt.Errorf("unknown flag: -%s", flagName[1])
		}

		return err
	}

	return nil
}

func newTabWriter(w io.Writer) *tabwriter.Writer {
	tw := tabwriter.NewWriter(w, 0, 0, 3, ' ', 0)

	return tw
}

func printMainHelp(w *tabwriter.Writer) {
	fmt.Fprintln(w, "Usage: lookctl <command> [arguments]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Commands:")
	fmt.Fprintln(w, "\tcurrent\tShow the currently used theme, icon, and cursor")
	fmt.Fprintln(w, "\tlist\tShow installed themes")
	fmt.Fprintln(w, "\tset\tSet the theme, icon, or cursor")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Run 'lookctl <command> -h' for more information on a command.")

	w.Flush()
}

func printListHelp(w *tabwriter.Writer) {
	fmt.Fprintln(w, "Usage: lookctl list [options]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "\t-cursor, --cursor\tShow installed cursor themes")
	fmt.Fprintln(w, "\t-gtk, --gtk\tShow installed themes (selected by default)")
	fmt.Fprintln(w, "\t-icon, --icon\tShow installed icon themes")

	w.Flush()
}

func printSetHelp(w *tabwriter.Writer) {
	fmt.Fprintln(w, "Usage: lookctl set [options] [arguments]")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "Options:")
	fmt.Fprintln(w, "\t-color-scheme, --color-scheme\tManually set cursor theme")
	fmt.Fprintln(w, "\t-cursor, --cursor\tSet cursor theme")
	fmt.Fprintln(w, "\t-gtk, --gtk\tSet theme")
	fmt.Fprintln(w, "\t-icon, --icon\tSet icon theme")

	w.Flush()
}

func printCurrentHelp(w *tabwriter.Writer) {
	fmt.Fprintln(w, "Usage: lookctl current")
	fmt.Fprintln(w, "Show applied themes")

	w.Flush()
}
