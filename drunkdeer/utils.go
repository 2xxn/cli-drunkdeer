package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/2xxn/cli-drunkdeer/driver"
	"github.com/fatih/color"
	"github.com/sstallion/go-hid"
)

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destination.Close()

	if _, err := io.Copy(destination, source); err != nil {
		return fmt.Errorf("failed to copy file contents: %w", err)
	}

	return destination.Sync()
}

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http") && strings.Contains(path, "://")
}

func getDevice(index int) (*hid.Device, error) {
	devices := FindDrunkDeerDevices()
	if len(devices) == 0 {
		return nil, fmt.Errorf("no devices found")
	}

	if index < 0 || index >= len(devices) {
		return nil, fmt.Errorf("invalid keyboard index: %d", index)
	}

	device, err := hid.OpenPath(devices[index].Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open device: %w", err)
	}

	return device, nil
}

func grabDeviceIdentity(deviceInfo *hid.DeviceInfo) (*driver.DDKeyboardIdentity, error) {
	device, err := hid.OpenPath(deviceInfo.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open device for identity: %w", err)
	}
	defer device.Close()

	controller := driver.NewDrunkDeerController(device)
	return controller.GetIdentity(), nil
}

func (a *App) getConfig(loadPath string) *Config {
	debugPrintf("Loading profile from %s", loadPath)

	var config Config
	var data []byte
	var err error

	if isURL(loadPath) {
		data, err = download(loadPath)
		handleError("Failed to download profile", err)
	} else {
		profilePath := a.resolveProfilePath(loadPath)
		data, err = os.ReadFile(profilePath)
		handleError("Failed to read profile file", err)
	}

	err = json.Unmarshal(data, &config)
	handleError("Failed to parse profile JSON", err)

	return &config
}

func (a *App) resolveProfilePath(loadPath string) string {
	if filepath.IsAbs(loadPath) {
		return loadPath
	}

	profilePath := filepath.Join(a.profilePath, loadPath)
	if !strings.HasSuffix(profilePath, ".json") {
		profilePath += ".json"
	}

	return profilePath
}

func handleError(message string, err error) {
	if err == nil {
		return
	}

	color.HiRed("%s: %v\n", message, err)
	os.Exit(1)
}
