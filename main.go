package main

// all imports
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"

	"./handlers"

	"github.com/gomidi/connect"
	driver "github.com/gomidi/rtmididrv"
	socketio "github.com/googollee/go-socket.io"

	//https://godoc.org/github.com/graarh/golang-socketio
	"github.com/zserge/webview"
)

var server socketio.Server

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

	type devicesInData struct {
		InDevice   int
		OutDevice  int
		DeviceType int
	}

	type deviceOutAddUI struct {
		DevName      string
		FriendlyName string
		Enabled      string
	}

	server.OnEvent("/", "AddDevice", func(s socketio.Conn, msg string) error {
		var data devicesInData
		err := json.Unmarshal([]byte(msg), &data)
		if err != nil {
			log.Println(err)
		}

		fmt.Println("devices: ", data.InDevice, data.InDevice)

		//check validity of the data
		if data.DeviceType == 0 {
			// it is a midi device therefore a listener is needed

			in, err := connect.OpenIn(drv, data.InDevice, "")

			//handles the potential error
			if err != nil {
				handlers.Must(err)
				if in.IsOpen() {
					in.Close()
				}
			}

			out, err := connect.OpenOut(drv, data.OutDevice, "")

			//handles the potential error
			if err != nil {
				handlers.Must(err)
				if out.IsOpen() {
					out.Close()
				}
			}

			//if the device is successfully opened it tries to attach a listener
			if in.IsOpen() {
				err := in.SetListener(handleMidiEvent)
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
		}

		//emit creation of the device to the UI
		fmt.Println("listening to device", data.InDevice)
		/* device := deviceOutAddUI{DevName: }
		s.Emit("AddDeviceReturn", json2text(device))
		*/
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
		URL:       "http://localhost/ui/",
		Resizable: true,
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

func cliLog(cause string, body string, threatLevel int) {
	server.sockets
}
