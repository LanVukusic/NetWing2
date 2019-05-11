package helpers

import (
	"encoding/json"
	"log"

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
		log.Println(err)
		return
	}

	switch cmd.Event {
	case "alert":
		Alert(w, cmd.Value)

	case "refresh_midi_devices":
		/* fmt.Println(drv.Ins()) */
		//fmt.Println("okey js klice nazaj kar je ql")
		midicode.SendMIDIdevices(w, midicode.GetMIdiDevices())
		//fmt.Println(getMIdiDevices(drv))
	}

}

// Alert creates an "alert" windows on the UI
func Alert(w webview.WebView, text string) {
	w.Eval("alert('" + text + "')")
}

// Must is an error handler
func Must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
