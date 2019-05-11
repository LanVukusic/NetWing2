package main

// all imports
import (
	"fmt"
	"os"
	"runtime"

	"./helpers"
	"./midicode"

	driver "github.com/gomidi/rtmididrv"
	"github.com/zserge/webview"
)

// main
func main() {
	runtime.LockOSThread()
	var err error
	midicode.Drv, err = driver.New()
	helpers.Must(err)

	fmt.Println("On start", midicode.GetMIdiDevices())

	// web view settings
	var rootDirectory, _ = os.Getwd()
	w := webview.New(webview.Settings{
		Width:                  800,
		Height:                 600,
		Title:                  "NetWing",
		URL:                    "file://" + rootDirectory + "/web/view/index.html",
		ExternalInvokeCallback: helpers.HandleRPC,
		Resizable:              true,
	})

	defer w.Exit()
	w.Run()
}
