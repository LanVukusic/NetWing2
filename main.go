package main

// all imports
import (
	"fmt"
	"net/http"
	"runtime"

	"./helpers"

	"github.com/zserge/webview"
)

//Counter is a bitch
type Counter struct {
	Value int `json:"value"`
}

// Add increases the value of a counter by n
func (c *Counter) Add(n int) {
	c.Value = c.Value + int(n)
}

// main
func main() {
	// start services
	fmt.Println("Starting webview")
	runtime.LockOSThread()

	//create a http server to serve the UI both remote and to the local client
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Testing")
		})

		fs := http.FileServer(http.Dir("web/view"))
		http.Handle("/ui/", http.StripPrefix("/ui/", fs))

		fs2 := http.FileServer(http.Dir("web/style"))
		http.Handle("/style/", http.StripPrefix("/style/", fs2))

		fs3 := http.FileServer(http.Dir("web/js"))
		http.Handle("/js/", http.StripPrefix("/js/", fs3))

		fs4 := http.FileServer(http.Dir("web/static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs4))

		http.ListenAndServe(":80", nil)
	}()

	// start midi driver
	/* var err error
	midicode.drv, err = driver.New()
	handlers.Must(err)
	defer midicode.drv.Close() */

	// juice up the OSC service
	/* fmt.Println("Starting OSC")
	osclib.StartOSCServer() */

	// web view settings

	wb := webview.New(webview.Settings{
		Width:  1400,
		Height: 800,
		Title:  "NetWing",
		/* URL:                    "file://" + rootDirectory + "/web/view/index.html", */
		URL:                    "http://localhost/ui/",
		ExternalInvokeCallback: helpers.HandleRPC,
		Resizable:              true,
	})

	defer wb.Exit()

	wb.Dispatch(func() {
		//Create necessary UI bindingsS
		wb.Bind("counter", &Counter{})
		wb.Eval("alert('asd')")
		wb.Eval("alert(counter)")
		/* updateMidiIns, _ := wb.Bind("MidiListIn", &MidiListIn{})
		updateMidiOuts, _ := wb.Bind("MidiListOut", &MidiListOut{}) */

	})

	wb.Run()

	/* // MidiDevice midi device interface
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
	} */

}

/* // Reset sets the value of a counter back to zero
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
} */

/* // Reset sets the value of a counter back to zero
func setval(c *MidiListOut) {
	c.Value = 0
} */

/* func updateMidiLists() {
	list := GetMIdiDevices()

}
*/
