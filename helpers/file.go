package helpers

import (
	"github.com/gomidi/connect"
	"github.com/zserge/webview"
)

// Alert creates an "alert" windows on the UI
func Alert(w webview.WebView, text string) {
	w.Eval("alert('" + text + "')")
}

//MidiDevice is a holder for a midi device that gets sent from BEnd to UI
type MidiDevice struct {
	NameWithID string
	Name       string
	ID         int
}

//MidiPackage holds in and ut devices for updating the UI
type MidiPackage struct {
	Event string
	Outs  []MidiDevice
	Ins   []MidiDevice
}

//InterfaceDevice holds a reference to a device that user assigns. It can be any type of interface that will potentially get mapped.
type InterfaceDevice struct {
	Active       bool
	BindID       int   // device ID for internal binding. Unique for every interface
	DeviceType   int16 // 0 = MIDI
	HardwareName string
	FriendlyName string
	HardwareID   int // let's say 0 as in first midi device
}

//InterfaceMessage is an incoming message from any device interface. It serves as the key in the interface - OSC mapping.
type InterfaceMessage struct {
	msgType int16 // sametype as InterfaceDevice: 0 = MIDI ...
	bindID  int   // same as InterfaceDevice bindID
	channel int32 // MAY HAVE TO CHANGE THAT
	value   int16 // the value normalized from 0 to 100
}

//OSCOutput has all the defined values to create an OSC message.
type OSCOutput struct {
	message  string
	argument int16
}

//CliMsg is a holder for any warn or error that gets sent from backend to be displayed on the CLI of clients
type CliMsg struct {
	Event       string
	Cause       string
	Body        string
	ThreatLevel int
}

//WSMsgTemplate represents the type of a message that gets returned from UI. It carries an event and some data.
type WSMsgTemplate struct {
	Event string `json:"Event"`
	Data  string `json:"Data"`
}

//Bind2MIDI will server as a dictionary to map MIDI devices and their bindIDs. Serves for reconnect purposes
type Bind2MIDI struct {
	BindID   int
	MidiPort connect.In
}

//GetOSCMessage returns an OSC message ready to get sent.
/* func (o OSCOOutput) GetOSCMessage() (msg string, err error) {

} */
