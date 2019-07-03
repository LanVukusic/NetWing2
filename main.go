package main

// all imports
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"./handlers"
	"./helpers"

	"github.com/gomidi/connect"
	driver "github.com/gomidi/rtmididrv"
	socketio "github.com/googollee/go-socket.io"
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

	// create Socket.IO server to handle comunication with frontend
	fmt.Println("Starting SocketIO connection")
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.LocalAddr())
		return nil
	})

	/* server.OnEvent("/", "startMidi", func(s socketio.Conn) error {
		fmt.Println("Starting MIDI service")

		return nil
	})

	server.OnEvent("/", "stopMidi", func(s socketio.Conn) error {
		fmt.Println("Stopping MIDI service")
		return nil
	}) */

	server.OnEvent("/", "refreshMidi", func(s socketio.Conn) error {
		fmt.Println("Refreshing device list")
		//generate ins and outs
		data, err := getMIDIDevices(drv)
		handlers.Must(err)

		//json-ify the data
		dataJ, err := json2text(data)
		handlers.Must(err)

		//emit data
		s.Emit("refreshMidiRet", dataJ)

		return nil
	})

	go server.Serve()
	defer server.Close()

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

		http.Handle("/socket.io/", server)

		http.ListenAndServe(":80", nil)
	}()

	// web view settings
	fmt.Println("Starting webview")
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
	wb.Run()
}

type midiDevice struct {
	Name string
	ID   int
}

type midiPackage struct {
	Outs []midiDevice
	Ins  []midiDevice
}

func getMIDIDevices(drv connect.Driver) (outDevices midiPackage, err error) {
	//gets the inputs
	ins, err := drv.Ins()
	if err != nil {
		return outDevices, err
	}
	for _, el := range ins {

		outDevices.Ins = append(outDevices.Ins, midiDevice{
			Name: el.String(),
			ID:   el.Number()})
	}

	outs, err := drv.Ins()
	if err != nil {
		return outDevices, err
	}
	for _, el := range outs {

		outDevices.Outs = append(outDevices.Outs, midiDevice{
			Name: el.String(),
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
