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
	BindID    int
	MidiPort  connect.In
	WasOnline bool
}

//MIDILearnMessage serves as an interface to comunicate to frontend what midi channel was to be used in bind to a fader / execWindow item
type MIDILearnMessage struct {
	Event     string `json:"Event"`
	Interf    int    `json:"Interf"`
	DeviceID  int    `json:"DeviceID"`
	ChannelID byte   `json:"ChannelID"`
}

//InternalDevice is used as a key in mappings HashMap. it tells us the interface, device, and channel
type InternalDevice struct {
	InterfaceType int // MIDI = 0
	DeviceID      int
	ChannelID     byte
}

//InternalOutput is used as a value in mappings hashmap. it tells what value and type of data should be processed
type InternalOutput struct {
	OutType float64
	OutChan int
	OutPage int
	Fade    bool
}

//MappingResponse is sent back to the client, to confirm that the interface is bound
type MappingResponse struct {
	Event     string `json:"Event"`
	Interface int    `json:"Interface"`
	DeviceID  int    `json:"DeviceID"`
	ChannelID byte   `json:"ChannelID"`
	FaderID   int    `json:"FaderID"`
	ExecID    int    `json:"ExecID"`
	ExecPage  int    `json:"ExecPage"`
}

//FaderUpdate is sent to client to update fader value with current MIDI value
type FaderUpdate struct {
	Event   string `json:"Event"`
	Type    int    `json:"Type"`
	FaderID int    `json:"FaderID"`
	Value   byte   `json:"Value"`
}

type ExecUpdate struct {
	Event    string `json:"Event"`
	Type     int    `json:"Type"`
	FaderID  int    `json:"FaderID"`
	PageID   int    `json:"PageID"`
	FadeType bool   `json:"FadeType"`
	Value    int    `json:"Value"`
}

// ExecWindow is a struct to be put in array, so server can keep up with active exec pages for saving and later updates
type ExecWindow struct {
	Event  string `json:"Event"`
	Page   int    `json:"Page"`
	Width  int    `json:"Width"`
	Height int    `json:"Height"`
}

//GetOSCMessage returns an OSC message ready to get sent.
/* func (o OSCOOutput) GetOSCMessage() (msg string, err error) {

} */
