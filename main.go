package main

// all imports
import (
	"fmt"
	"os"
	"runtime"

	"./helpers"

	"github.com/gomidi/connect"
	"github.com/zserge/webview"
)

// main
func main() {
	// start services
	fmt.Println("Starting webview")
	runtime.LockOSThread()

	// start midi driver
	/* var err error
	midicode.drv, err = driver.New()
	handlers.Must(err)
	defer midicode.drv.Close() */

	// juice up the OSC service
	/* fmt.Println("Starting OSC")
	osclib.StartOSCServer() */

	// web view settings
	var rootDirectory, _ = os.Getwd()
	wb := webview.New(webview.Settings{
		Width:                  800,
		Height:                 600,
		Title:                  "NetWing",
		URL:                    "file://" + rootDirectory + "/web/view/index.html",
		ExternalInvokeCallback: helpers.HandleRPC,
		Resizable:              true,
	})

	defer wb.Exit()
	wb.Run()
	// MidiDevice midi device interface
	type MidiDevice struct {
		Name string `json:"name"`
		ID   int    `json:"id"`
	}

	// MidiInList list of In MIDI devices
	type MidiListIn struct {
		Devices []MidiDevice `json:"devices"`
	}
	type MidiListOut struct {
		Devices []MidiDevice `json:"devices"`
	}

	wb.Dispatch(func() {
		//Create necessary UI bindings
		wb.Bind("counter", &Counter{})
		updateMidiIns, _ := wb.Bind("MidiListIn", &MidiListIn{})
		updateMidiOuts, _ := wb.Bind("MidiListOut", &MidiListOut{})
		fmt.Println(updateMidiIns, updateMidiOuts)
	})

}

type Counter struct {
	Value int `json:"value"`
}

// Add increases the value of a counter by n
func (c *Counter) Add(n int) {
	c.Value = c.Value + int(n)
}

// Reset sets the value of a counter back to zero
func (c *Counter) Reset() {
	c.Value = 0
}

//Midi section
var drv connect.Driver

// GetMIdiDevices gets the list of available devices from the OS
func GetMIdiDevices() (outs []connect.In) {
	ins, _ := drv.Ins()

	for _, el := range ins {
		outs = append(outs, el)
	}
	return outs
}

/* // Reset sets the value of a counter back to zero
func setval(c *MidiListOut) {
	c.Value = 0
} */

/* func updateMidiLists() {
	list := GetMIdiDevices()

}
*/
