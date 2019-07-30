package helpers

import (
	"github.com/zserge/webview"
)

// Alert creates an "alert" windows on the UI
func Alert(w webview.WebView, text string) {
	w.Eval("alert('" + text + "')")
}

//MidiDevice is a holder for a midi device that gets sent from BEnd to UI
type MidiDevice struct {
	Name string
	ID   int
}

//MidiPackage holds in and ut devices for updating the UI
type MidiPackage struct {
	Outs []MidiDevice
	Ins  []MidiDevice
}

//InterfaceDevice holds a reference to a device that user assigns. It can be any type of interface that will potentially get mapped.
type InterfaceDevice struct {
	BindID       int16
	DeviceType   int16 // 0 = MIDI
	HardwareName string
	FriendlyName string
	HardwareID   int
}

//InterfaceMessage is an incoming message from any device interface. It serves as the key in the interface - OSC mapping.
type InterfaceMessage struct {
	msgType int16 // sametype as InterfaceDevice: 0 = MIDI ...
	bindID  int16 // same as InterfaceDevice bindID
	channel int32 // MAY HAVE TO CHANGE THAT
	value   int16 // the value normalized from 0 to 100
}

//OSCoutput has all the defined values to create an OSC message.
type OSCoutput struct {
	message  string
	argument int16
}

//GetOSCMessage returns an OSC message ready to get sent.
/* func (o OSCoutput) GetOSCMessage() (msg string, err error) {

} */
