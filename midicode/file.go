package midicode

import (
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
