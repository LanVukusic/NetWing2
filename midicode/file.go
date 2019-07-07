package midicode

import (
	"fmt"
	"strconv"

	"github.com/gomidi/connect"
	"github.com/zserge/webview"
)

// Drv is gay
var Drv connect.Driver

// SendMIDIdevices sends the list of devices on the UI and displays them.
func SendMIDIdevices(w webview.WebView, devices []connect.In) {

	for _, el := range devices {
		//fmt.Println("addDevice ('" + el.String() + "', " + strconv.Itoa(el.Number()) + ")")
		w.Eval("addDevice ('" + el.String() + "', " + strconv.Itoa(el.Number()) + ")")
	}
}

// GetMIdiDevices gets the list of available devices from the OS
func GetMIdiDevices() (outs []connect.In) {
	ins, _ := Drv.Ins()

	for _, el := range ins {
		outs = append(outs, el)
	}
	return outs
}

/* //ListenMidi checks availibility and ATTACHES a MIDI listener to the device
func ListenMidi(id int) {

	// this line gets the device "id" from the driver, opens it and returns the active device
	in, err := connect.OpenIn(Drv, id, "")

	//handles the potential error
	if err != nil {
		handlers.Must(err)
		if in.IsOpen() {
			in.Close()
		}
	}

	//if the device is successfully opened it tries to attach a listener
	if in.IsOpen() {
		err := in.SetListener(handleMIDIevent)
		//if unsuccessful, handle the error
		if err != nil {
			handlers.Must(err)
			if in.IsOpen() {
				//stop the device
				in.StopListening()
				in.Close()
			}
		}
	}
} */

//StopListenMidi checks availability and DETACHES a MIDI listener from the device
func StopListenMidi(id int) {
	//gets the devices from driver
	ins, _ := Drv.Ins()
	//stops the listener
	ins[id].StopListening()
	//checks if the device is open and if so, it closes it
	if ins[id].IsOpen() {
		ins[id].Close()
	}
}

func handleMIDIevent(data []byte, deltaMicroseconds int64) {
	//fmt.Println(data)
	//osclib.SendOSC(int(data[1]), int(data[2]))
	msg := "Chan: " + strconv.Itoa(int(data[1])) + " Value: " + strconv.Itoa((int(data[2])))
	//handlers.W.Eval("cliLog (0, MIDI in, " + msg + ");")
	fmt.Println(msg)
}
