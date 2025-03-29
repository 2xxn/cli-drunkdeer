package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func (a *App) handleArgs() {
	if a.args.Debug {
		debug = true
		debugPrintf("DEBUG ENABLED")
	}

	switch a.args.Command {
	case "load":
		a.args.Load = a.args.CmdValue
	case "save":
		a.args.Save = a.args.CmdValue
	case "reset":
		a.args.Reset = true
	case "import":
		a.args.Import = a.args.CmdValue
	case "version":
		a.showVersion()
	case "list":
		displayDeviceList()
	case "profiles":
		a.displayProfiles()
	}

	switch {
	case a.args.Version:
		a.showVersion()
	case a.args.List:
		displayDeviceList()
	case a.args.Profiles:
		a.displayProfiles()
	case a.args.Import != "":
		a.importConfig(a.args.Import)
	case a.args.Save != "":
		a.saveConfig(a.args.Save)
	}

	a.keyboardIndex = a.args.Index
	if a.keyboardIndex < 0 {
		a.keyboardIndex = 0
	}
}

func (a *App) showVersion() {
	fmt.Println("Version 1.0.0")
	os.Exit(0)
}

func (a *App) saveConfig(profilePath string) {
	if isURL(profilePath) {
		a.saveConfigFromURL(profilePath)
		return
	}

	if !filepath.IsAbs(profilePath) {
		color.HiRed("Error: %s is not an absolute path\n", profilePath)
		os.Exit(1)
	}

	fileName := ensureJSONExtension(filepath.Base(profilePath))
	targetPath := filepath.Join(a.profilePath, fileName)

	if _, err := os.Stat(targetPath); err == nil {
		color.HiRed("Error: %s already exists\n", targetPath)
		os.Exit(1)
	}

	if err := copyFile(profilePath, targetPath); err != nil {
		color.HiRed("Error saving profile: %v\n", err)
		os.Exit(1)
	}

	color.HiGreen("Profile saved to %s\n", targetPath)
	os.Exit(0)
}

func (a *App) saveConfigFromURL(url string) {
	data, err := download(url)
	if err != nil {
		color.HiRed("Error downloading profile: %v\n", err)
		os.Exit(1)
	}

	fileName := ensureJSONExtension(
		strings.TrimSuffix(
			filepath.Base(url),
			filepath.Ext(url),
		),
	)
	targetPath := filepath.Join(a.profilePath, fileName)

	file, err := os.Create(targetPath)
	if err != nil {
		color.HiRed("Error creating profile file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	if _, err := io.WriteString(file, string(data)); err != nil {
		color.HiRed("Error writing profile: %v\n", err)
		os.Exit(1)
	}

	color.HiGreen("Profile saved to %s\n", targetPath)
	os.Exit(0)
}

func (a *App) importConfig(source string) {
	var data []byte
	var err error

	if isURL(source) {
		data, err = download(source)
	} else {
		data, err = os.ReadFile(source)
	}

	if err != nil {
		color.HiRed("Error loading profile: %v\n", err)
		os.Exit(1)
	}

	config := parseDrunkDeerConfig(data).convertToCLIConfig()
	fileName := ensureJSONExtension(filepath.Base(source))
	targetPath := filepath.Join(a.profilePath, fileName)

	if err := a.writeConfigToFile(config, targetPath); err != nil {
		color.HiRed("Error saving imported profile: %v\n", err)
		os.Exit(1)
	}

	a.showImportSuccess(fileName, targetPath)
}

func (a *App) writeConfigToFile(config *Config, path string) error {
	jsonData, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(path, jsonData, 0644)
}

func (a *App) showImportSuccess(fileName, path string) {
	profileName := strings.TrimSuffix(fileName, ".json")
	color.HiGreen("Profile imported to %s\n", path)
	fmt.Printf("Use: ")
	color.HiBlue("drunkdeer load %s\n", profileName)
	os.Exit(0)
}

func (a *App) displayProfiles() {
	profiles, err := os.ReadDir(a.profilePath)
	if err != nil {
		color.HiRed("Error reading profile directory: %v", err)
		os.Exit(1)
	}

	if len(profiles) == 0 {
		color.HiRed("No profiles found")
		os.Exit(0)
	}

	color.HiGreen("Available profiles:")
	for _, profile := range profiles {
		if profile.IsDir() {
			continue
		}
		profileName := strings.TrimSuffix(profile.Name(), ".json")
		color.White("%s", profileName)
	}
	os.Exit(0)
}

func displayDeviceList() {
	devices := FindDrunkDeerDevices()
	if len(devices) == 0 {
		color.HiRed("No devices found")
		os.Exit(0)
	}

	color.HiGreen("Connected devices:")
	for i, device := range devices {
		identity, err := grabDeviceIdentity(&device)
		if err != nil {
			color.HiRed("Error reading identity for device %d: %v\n", i, err)
			continue
		}

		firmware := color.RGB(0x80, 0x80, 0x80).Sprintf(
			"(firmware version: v%v)",
			identity.FirmwareVersion,
		)
		fmt.Printf("%v: %v %v\n",
			color.WhiteString("%d", i),
			color.HiBlueString("DrunkDeer "+identity.KeyboardModel),
			firmware,
		)
	}
	os.Exit(0)
}

func ensureJSONExtension(filename string) string {
	if !strings.HasSuffix(filename, ".json") {
		return filename + ".json"
	}
	return filename
}
