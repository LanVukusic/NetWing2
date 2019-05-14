package main

// all imports
import (
	"fmt"
	"os"
	"runtime"

	"./handlers"
	"./helpers"
	"./midicode"
	"./osclib"

	driver "github.com/gomidi/rtmididrv"
	"github.com/zserge/webview"
)

// main
func main() {
	fmt.Println("Starting webview")
	runtime.LockOSThread()
	var err error
	midicode.Drv, err = driver.New()
	handlers.Must(err)
	defer midicode.Drv.Close()

	fmt.Println("Starting OSC")
	osclib.StartOSCserver()

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
