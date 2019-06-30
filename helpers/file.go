package helpers

import (
	"encoding/json"
	"strconv"

	"../handlers"
	"../midicode"
	"github.com/zserge/webview"
)

// HandleRPC this handles the webview call
func HandleRPC(w webview.WebView, data string) {
	cmd := struct {
		Event string `json:"type"`
		Value string `json:"value"`
	}{}
	if err := json.Unmarshal([]byte(data), &cmd); err != nil {
		handlers.Must(err)
		return
	}

	switch cmd.Event {
	case "alert":
		Alert(w, cmd.Value)

	case "refresh_midi_devices":
		midicode.SendMIDIdevices(w, midicode.GetMIdiDevices())

	case "clear_midi_devices":
		w.Eval("clearDevices();")

	case "listen_debug_midi_devices":
		i, err := strconv.Atoi(cmd.Value)
		handlers.Must(err)
		midicode.ListenMidi(i)
	}

}

// Alert creates an "alert" windows on the UI
func Alert(w webview.WebView, text string) {
	w.Eval("alert('" + text + "')")
}




