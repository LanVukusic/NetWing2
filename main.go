package main

// all imports
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"

	"./helpers"

	"github.com/gomidi/connect"
	driver "github.com/gomidi/rtmididrv"
	"github.com/gorilla/websocket"

	"github.com/zserge/webview"
)

// main mapping dictionary
var mainmappings map[helpers.InterfaceMessage]int

// main device array. supports max of 50 devices
var mainDeviceList [50]helpers.InterfaceDevice

var drvMIDI connect.Driver

var upgrader websocket.Upgrader
var wsConnections []*websocket.Conn // supports max 30 clients

// main
func main() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	mainmappings = make(map[helpers.InterfaceMessage]int)

	// start services
	fmt.Println("locking Threads")
	runtime.LockOSThread()

	//init midi
	drvMIDI = nil

	/* //start OSC
	fmt.Println("Starting OSC")
	osclib.StartOSCServer() */

	// create a http server to serve the UI both remote and to the local client
	fmt.Println("Starting webserver")
	go runWebserver()

	// web view settings

	if true {
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

	cliLog("Engine", "Engine running smoothly", 0)

}

func handleWSMessage(messageType int, p []byte, socket *websocket.Conn) {
	fmt.Println(string(p), socket.LocalAddr().String())

	/* if err := socket.WriteMessage(messageType, p); err != nil {
		handleErr(err, "Websocket message writing error for: "+socket.LocalAddr().String(), true)
	} */
	//cliLog("testing", "hey testing cli working", 0)
}

func broadcastMessage(msg interface{}) {
	for _, client := range wsConnections {
		err := client.WriteJSON(msg)
		if err != nil {
			handleErr(err, "error during websocket broadcasting ", false)
		}
	}
}

func upgradeConnection(w http.ResponseWriter, r *http.Request) {
	// upgrade connection to WS connection
	wsConnection, err := upgrader.Upgrade(w, r, nil)
	wsConnections = append(wsConnections, wsConnection)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("websocket connection from: " + wsConnection.LocalAddr().String())
	cliLog("WebServer", fmt.Sprintf("Client Connected : %v", wsConnection.LocalAddr().String()), 0)

	go func() {
		for {
			messageType, p, err := wsConnection.ReadMessage()
			if err != nil {
				handleErr(err, "Websocket message reading error for: "+wsConnection.LocalAddr().String(), true)
				return
			}
			handleWSMessage(messageType, p, wsConnection)
		}
	}()
}

func runWebserver() {
	// handle root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Testing")
	})

	// handle web socket connection
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		upgradeConnection(w, r)
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
}

func getMIDIDevices(drv connect.Driver) (outDevices helpers.MidiPackage, err error) {
	//gets the inputs
	ins, err := drv.Ins()
	if err != nil {
		return outDevices, err
	}
	for _, el := range ins {

		outDevices.Ins = append(outDevices.Ins, helpers.MidiDevice{
			// it splits the string by spaces, removeslast slice (the number) and joins it back together to form a string
			Name: strings.Join(strings.Fields(el.String())[:len(strings.Fields(el.String()))-1], ""),
			ID:   el.Number()})
	}

	outs, err := drv.Ins()
	if err != nil {
		return outDevices, err
	}
	for _, el := range outs {

		outDevices.Outs = append(outDevices.Outs, helpers.MidiDevice{
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
	alert := helpers.CliMsg{
		Event:       "cli",
		Cause:       cause,
		Body:        body,
		ThreatLevel: threatLevel,
	}
	broadcastMessage(alert)
}

func handleErr(err error, msg string, bcast bool) {
	fmt.Println(err, msg)
	if bcast {
		alert := helpers.CliMsg{
			Cause:       err.Error(),
			Body:        msg,
			ThreatLevel: 2,
		}

		broadcastMessage(alert)
	}

}

//ListenMidi checks availability and ATTACHES a MIDI listener to the device
func ListenMidi(id int) {
	// this line gets the device "id" from the driver, opens it and returns the active device
	in, err := connect.OpenIn(drvMIDI, id, "")

	//handles the potential error
	if err != nil {
		handleErr(err, "MIDI device:"+string(id)+"is unavailble", true)
		if in.IsOpen() {
			in.Close()
		}
	}
	//if the device is successfully opened it tries to attach a listener
	if in.IsOpen() {
		err := in.SetListener(handleMidiEvent)
		//if unsuccessful, handle the error
		if err != nil {
			handleErr(err, "MIDI device:"+string(id)+"cant receive a listener", true)
			if in.IsOpen() {
				//stop the device
				in.StopListening()
				in.Close()
			}
		}
	}
	cliLog("MIDI Interf.", fmt.Sprintf("Listening to MIDI device %v", id), 0)
}

//StopListenMidi checks availability and DETACHES a MIDI listener from the device
func StopListenMidi(id int) {
	//gets the devices from driver
	ins, _ := drvMIDI.Ins()
	//stops the listener
	ins[id].StopListening()
	//checks if the device is open and if so, it closes it
	if ins[id].IsOpen() {
		ins[id].Close()
	}
	cliLog("MIDI", fmt.Sprintf("MIDI device %v successfully closed", id), 0)
}

func handleMIDIevent(data []byte, deltaMicroseconds int64) {
	//fmt.Println(data)
	//osclib.SendOSC(int(data[1]), int(data[2]))
	msg := "Chan: " + strconv.Itoa(int(data[1])) + " Value: " + strconv.Itoa((int(data[2])))
	//handlers.W.Eval("cliLog (0, MIDI in, " + msg + ");")
	fmt.Println(msg)
}

// function loops through devices and initializes their drivers and updates the UI
func initializeSavedDevices(devlist []helpers.InterfaceDevice) (err error) {
	for _, device := range devlist {
		//check for type and create the driver
		if device.DeviceType == 0 { // it's a MIDI device

			//if MIDI has no active driver create one
			if drvMIDI == nil {
				var err error
				drvMIDI, err = driver.New()
				if err != nil {
					handleErr(err, "cant create MIDI driver", true)
					return err
				}
			}

			//The driver is active, so assign a listener to the MIDI device.
			ListenMidi(device.HardwareID)

		} else {
			handleErr(nil, "unrecognized device type", true)
		}

		//update ui

		//check if device exists

	}
	cliLog("Initialization", "Init completed. Devices loaded successfully", 0)
	return nil
}
