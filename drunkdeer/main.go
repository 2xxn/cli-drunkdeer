package main

import (
	"time"

	"github.com/sstallion/go-hid"
)

const (
	defaultProfilePath        = "~/.drunkdeer"
	defaultWaitPerInstruction = 100 * time.Millisecond
)

var debug = false

func main() {
	hid.Init()
	defer hid.Exit()

	app := &App{profilePath: defaultProfilePath}

	app.parseArgs()
	app.setupProfilePath()
	app.handleArgs()
	app.setupDevice()

	defer app.cleanup()

	app.run()
}
