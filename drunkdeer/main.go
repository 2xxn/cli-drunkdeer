package main

import (
	"fmt"
	"os"
	"time"

	"github.com/2xxn/cli-drunkdeer/driver"
	"github.com/alexflint/go-arg"
	"github.com/fatih/color"
	"github.com/sstallion/go-hid"
)

const (
	defaultKeyboardIndex      = 0
	defaultProfilePath        = "~/.drunkdeer"
	defaultWaitPerInstruction = 100 * time.Millisecond
)

var (
	debug = false
)

type App struct {
	keyboardIndex int
	device        *hid.Device
	controller    *driver.DrunkDeerController
	profilePath   string
	args          Args
}

func debugPrintf(format string, v ...interface{}) {
	if !debug {
		return
	}

	if format[len(format)-1] != '\n' {
		format += "\n"
	}
	debugPrefix := color.HiGreenString("[DEBUG] ")
	fmt.Printf(debugPrefix+format, v...)
}

func main() {
	hid.Init()
	defer hid.Exit()

	app := &App{
		keyboardIndex: defaultKeyboardIndex,
		profilePath:   defaultProfilePath,
	}

	app.parseArgs()
	app.setupProfilePath()
	app.handleArgs()
	app.setupDevice()
	defer app.cleanup()

	app.run()
}

func (a *App) parseArgs() {
	arg.MustParse(&a.args)
	debug = a.args.Debug
}

func (a *App) setupProfilePath() {
	if a.profilePath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		handleError("Error getting home directory", err)

		a.profilePath = homeDir + a.profilePath[1:]
	}

	if err := os.MkdirAll(a.profilePath, 0755); err != nil {
		handleError("Error creating profile directory", err)
	}
}

func (a *App) setupDevice() {
	var err error
	a.device, err = getDevice(a.keyboardIndex)
	handleError("Error:", err)
	debugPrintf("Device opened")

	a.controller = driver.NewDrunkDeerController(a.device)
	a.controller.GetIdentity()
	debugPrintf("Created controller")

	if debug {
		a.controller.SetDebug(true)
	}
}

func (a *App) cleanup() {
	if a.device != nil {
		a.device.Close()
	}
}

func (a *App) run() {
	switch {
	case a.args.Reset:
		a.handleReset()
	case a.args.Load != "":
		a.handleLoadProfile()
	default:
		a.showHelp()
	}
}

func (a *App) handleReset() {
	color.HiRed("Resetting device to default settings")
	a.controller.WriteDefaults()
	time.Sleep(10 * defaultWaitPerInstruction)
	color.White("Reset complete")
	time.Sleep(defaultWaitPerInstruction)
}

func (a *App) handleLoadProfile() {
	config := a.getConfig(a.args.Load)
	if config.Model != "" && config.Model != a.controller.GetIdentity().KeyboardModel {
		color.HiRed("Profile model does not match device model (expected %s, got %s)",
			a.controller.GetIdentity().KeyboardModel, config.Model)
		os.Exit(1)
	}

	debugPrintf("Model: %v | Turbo: %v | RT: %v | Default actuation: %v",
		config.Model, config.Turbo, config.RapidTrigger.Enabled, config.DefaultActuation)

	actuations, downstrokes, upstrokes := a.prepareKeySettings(config)

	if !config.Light.Enabled {
		config.Light.Sequence = driver.SEQUENCE_OFF
	}

	a.configureLights(config)
	a.applySettings(config, actuations, downstrokes, upstrokes)

	color.White("Loaded %s%s%s",
		color.GreenString(a.args.Load),
		color.WhiteString(" for "),
		color.HiBlueString("DrunkDeer %s", config.Model))
	debugPrintf("Profile loaded")
	time.Sleep(defaultWaitPerInstruction)
}

func (a *App) prepareKeySettings(config *Config) ([]byte, []byte, []byte) {
	actuations := make([]byte, len(driver.KEYBOARD_LAYOUT))
	downstrokes := make([]byte, len(driver.KEYBOARD_LAYOUT))
	upstrokes := make([]byte, len(driver.KEYBOARD_LAYOUT))

	defaultAct := driver.ActuationFloatToByte(config.DefaultActuation)
	defaultDS := driver.ActuationFloatToByte(config.RapidTrigger.DefaultDownstroke)
	defaultUS := driver.ActuationFloatToByte(config.RapidTrigger.DefaultUpstroke)

	for i := range actuations {
		actuations[i] = defaultAct
		downstrokes[i] = defaultDS
		upstrokes[i] = defaultUS
	}

	for key, value := range config.ActuationPoints {
		i := driver.GetIndexByKey(key)
		actuations[i] = driver.ActuationFloatToByte(value)
	}

	for key, value := range config.RapidTriggers {
		i := driver.GetIndexByKey(key)
		downstrokes[i] = driver.ActuationFloatToByte(value[0])
		upstrokes[i] = driver.ActuationFloatToByte(value[1])
	}

	return actuations, downstrokes, upstrokes
}

func (a *App) configureLights(config *Config) {
	a.controller.Light = &driver.DDLight{
		Sequence:   byte(config.Light.Sequence),
		Speed:      byte(config.Light.Speed),
		Direction:  byte(config.Light.Direction),
		Brightness: byte(config.Light.Brightness),
	}
}

func (a *App) applySettings(config *Config, actuations, downstrokes, upstrokes []byte) {
	a.controller.SendRapidTriggerTurbo(config.RapidTrigger.Enabled, config.Turbo)
	a.controller.SendLEDModeSelect(
		a.controller.Light.Direction,
		a.controller.Light.Sequence,
		a.controller.Light.Speed,
		a.controller.Light.Brightness,
		0xff,
	)
	a.controller.LoadActuations(actuations)
	a.controller.LoadDownstrokes(downstrokes)
	a.controller.LoadUpstrokes(upstrokes)
	time.Sleep(10 * defaultWaitPerInstruction)
}

func (a *App) showHelp() {
	if a.args.Command != "" {
		color.HiRed("Unknown command: %s\n", color.HiWhiteString(a.args.Command))
	}

	color.HiBlue("List of commands")
	color.White("For descriptions, run: drunkdeer --help")
	color.HiWhite("  - drunkdeer import <url/path>")
	color.HiWhite("  - drunkdeer load <profile>")
	color.HiWhite("  - drunkdeer save <profile>")
	color.HiWhite("  - drunkdeer profiles")
	color.HiWhite("  - drunkdeer reset")
	color.HiWhite("  - drunkdeer list")
	color.HiWhite("  - drunkdeer version")
}
