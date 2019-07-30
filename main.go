package main

// all imports
import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	"./handlers"
	"./helpers"

	"github.com/gomidi/connect"
	driver "github.com/gomidi/rtmididrv"

	"github.com/zserge/webview"
)

// main
func main() {
	// start services
	fmt.Println("locking Threads")
	runtime.LockOSThread()

	//start midi service
	fmt.Println("Starting MIDI service")
	drv, err := driver.New()
	handlers.Must(err)
	defer drv.Close()

	/* //start OSC
	fmt.Println("Starting OSC")
	osclib.StartOSCServer() */

	// create a http server to serve the UI both remote and to the local client
	fmt.Println("Starting webserver")
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

	// web view settings
	fmt.Println("Starting webview")
	wb := webview.New(webview.Settings{
		Width:  1400,
		Height: 800,
		Title:  "NetWing",
		/* URL:                    "file://" + rootDirectory + "/web/view/index.html", */
		URL:       "http://localhost/ui/",
		Resizable: true,
	})

	defer wb.Exit()
	wb.Run()
}

func getMIDIDevices(drv connect.Driver) (outDevices helpers.MidiPackage, err error) {
	//gets the inputs
	ins, err := drv.Ins()
	if err != nil {
		return outDevices, err
	}
	for _, el := range ins {

		outDevices.Ins = append(outDevices.Ins, midiDevice{
			// it splits the string by spaces, removeslast slice (the number) and joins it back together to form a string
			Name: strings.Join(strings.Fields(el.String())[:len(strings.Fields(el.String()))-1], ""),
			ID:   el.Number()})
	}

	outs, err := drv.Ins()
	if err != nil {
		return outDevices, err
	}
	for _, el := range outs {

		outDevices.Outs = append(outDevices.Outs, midiDevice{
			Name: strings.Join(strings.Fields(el.String())[:len(strings.Fields(el.String()))-1], ""),
			ID:   el.Number()})
	}
	return outDevices, nil
}

func json2text(in interface{}) (out string, err error) {
	var jsonData []byte
	jsonData, err = json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func handleMidiEvent(in []byte, time int64, deviceID int) {
	fmt.Println(fmt.Sprintf("Chn: %s, Val: %s, Device: %s", in[1], in[2], deviceID))
}

/* func cliLog(cause string, body string, threatLevel int) {
	server.sockets
} */
