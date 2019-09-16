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
	"time"

	"./helpers"

	"github.com/gobuffalo/packr"
	"github.com/gomidi/connect"
	driver "github.com/gomidi/rtmididrv"
	"github.com/gorilla/websocket"

	"github.com/zserge/webview"
)

// packeging
var boxView packr.Box
var boxStyle packr.Box
var boxStatic packr.Box
var boxJs packr.Box

// main mapping dictionary
var mainmappings map[helpers.InterfaceMessage]int

// main device array. supports max of 50 devices
var mainDeviceList [50]helpers.InterfaceDevice

// mapped interfaces and their bind ids.
var midi2idMappings [50]helpers.Bind2MIDI

// MIDI stuff
var drvMIDI connect.Driver

// Websocket and network
var upgrader websocket.Upgrader
var wsConnections []*websocket.Conn // supports max 30 clients

// main
func main() {
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// packages all static files in one binary
	boxView = packr.NewBox("./web/view")
	boxStatic = packr.NewBox("./web/")

	mainmappings = make(map[helpers.InterfaceMessage]int)

	// start services
	fmt.Println("locking Threads")
	runtime.LockOSThread()

	//start loopcheck to alert for disconnected devices
	go doEvery(2000*time.Millisecond, loopCheck)

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
		cliLog("Engine", "Engine running GUI mode", 0)
		defer wb.Exit()
		wb.Run()
	} else {
		cliLog("Engine", "Engine running CLI mode", 0)
		go func() {
			for true {
				// run infinite loop so the thread does not terminate
				//stupid me ofc it terminates.... it runs in a different thread
			}
		}()
	}

}

func midiDisconnected() {

}

// helper functions and whatnot
func loopCheck() {
	devices, _ := getMIDIDevices(drvMIDI)

	// check for inputs
	ins := devices.Ins

	for i, addedIterf := range midi2idMappings {
		if addedIterf.MidiPort == nil { // check if we are out of interfaces
			if i == 0 {
				//fmt.Println("works")
				return
			}
		}
		for _, allIterf := range ins {
			if addedIterf.MidiPort.String() == allIterf.NameWithID {
				fmt.Println("works", addedIterf.BindID)
				break
			} else {
				fmt.Println("works not")
			}
		}
	}
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func addUIdevice(deviceType int16, hName string, fName string, devID int, socket *websocket.Conn) {
	data := helpers.WSMsgTemplate{
		Event: "UiAddDevice",
		Data:  "{'ID':'" + strconv.Itoa(devID) + "','Hname':'" + hName + "','FriendlyName':'" + fName + "'}",
	}
	socket.WriteJSON(data)
}

func handleWSMessage(messageType int, p []byte, socket *websocket.Conn) {
	var raw map[string]interface{}
	err := json.Unmarshal(p, &raw)
	if err != nil {
		handleErr(err, fmt.Sprintf("Error while parsing JSON from %s", socket.LocalAddr().String()), true)
	}

	switch raw["event"] {
	case "getMidiDevices":
		data, err := getMIDIDevices(drvMIDI)
		if err != nil {
			handleErr(err, "Error while getting MIDI devices", true)
		} else {
			socket.WriteJSON(data)
		}
		break
	case "addInterface":
		devType := int(raw["deviceType"].(float64))

		if devType == 0 { // it is a MIDI device
			var id int
			id = 100
			// add device to main device list
			for i, interf := range mainDeviceList {
				if !interf.Active {
					id = i
					// free space in our device array
					mainDeviceList[i] = helpers.InterfaceDevice{
						Active:       true,
						BindID:       id, // that's the id for the device on the backend. not to be confused with midi ID
						DeviceType:   0,
						HardwareName: fmt.Sprintf("%v", raw["HardwareName"]),
						FriendlyName: fmt.Sprintf("%v", raw["FriendlyName"])}
					break
				} else {
					if interf.HardwareName == raw["HardwareName"] {
						// device allready exists
						cliLog("MIDI", "Midi device already active", 1)
						return
						break

					}
				}
			}

			if id > 50 {
				cliLog("Device level reached", "You have too many active devices.", 1)
				return
			}

			// handle its inputs / outputs
			ListenMidi(int(raw["inDevice"].(float64)), id)
			hName := fmt.Sprintf("%v", raw["HardwareName"])
			fName := fmt.Sprintf("%v", raw["FriendlyName"])
			sID := strconv.Itoa(id)

			fmt.Println(sID)
			data := helpers.WSMsgTemplate{
				Event: "UiAddDevice",
				Data:  "{'ID':'" + sID + "','Hname':'" + hName + "','FriendlyName':'" + fName + "'}",
			}

			// update the UI
			broadcastMessage(data)

		}
	}

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

	initializeUi(mainDeviceList[:], wsConnection)
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

	//serve index
	index, err := boxView.FindString("index.html")
	if err != nil {
		handleErr(err, "cant serve index", false)
	}
	http.HandleFunc("/ui/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, index)
	})

	/* http.Handle("/style/", http.StripPrefix("/style/", http.FileServer(boxStyle))) */
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(boxStatic)))
	/* http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(boxJs))) */
	http.ListenAndServe(":80", nil)
}

func getMIDIDevices(drv connect.Driver) (outDevices helpers.MidiPackage, err error) {
	outDevices.Event = "refreshMidiRet"
	if drv == nil {
		drv, err = driver.New()
		if err != nil {
			handleErr(err, "Failed creation of MIDI device driver.", true)
			return
		}
	}

	//gets the inputs
	ins, err := drv.Ins()
	if err != nil {
		return outDevices, err
	}

	for _, el := range ins {

		outDevices.Ins = append(outDevices.Ins, helpers.MidiDevice{
			// it splits the string by spaces, removeslast slice (the number) and joins it back together to form a string
			NameWithID: el.String(),
			Name:       strings.Join(strings.Fields(el.String())[:len(strings.Fields(el.String()))-1], ""),
			ID:         el.Number()})
	}

	outs, err := drv.Outs()
	if err != nil {
		return outDevices, err
	}
	for _, el := range outs {

		outDevices.Outs = append(outDevices.Outs, helpers.MidiDevice{
			NameWithID: el.String(),
			Name:       strings.Join(strings.Fields(el.String())[:len(strings.Fields(el.String()))-1], ""),
			ID:         el.Number()})
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
	fmt.Println(fmt.Sprintf("Chn: %s, Val: %s, Device: %s", int(in[1]), int(in[2]), int(deviceID)))
	cliLog("MIDI", fmt.Sprintf("Chn: %v, Val: %v, Device: %v", int(in[1]), int(in[2]), int(deviceID)), 0)
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
		cliLog(err.Error(), msg, 2)
	}

}

//ListenMidi checks availability and ATTACHES a MIDI listener to the device
func ListenMidi(id int, bind int) (err error) {
	// check if the MIDI driver is missing
	if drvMIDI == nil {
		var err error
		drvMIDI, err = driver.New()
		if err != nil {
			handleErr(err, "Error creating midi driver", true)
			return err
		}
	}

	// this line gets the device "id" from the driver, opens it and returns the active device
	in, err := connect.OpenIn(drvMIDI, id, "")

	//handles the potential error
	if err != nil {
		handleErr(err, "MIDI device:"+string(id)+"is unavailable", true)
		if in.IsOpen() {
			in.Close()
		}
		return err
	}

	//if the device is successfully opened it tries to attach a listener
	if in.IsOpen() {

		//add a binding to the monitoring array.
		for i, interf := range midi2idMappings {
			if interf.MidiPort == nil {
				midi2idMappings[i] = helpers.Bind2MIDI{
					BindID:   bind,
					MidiPort: in}
				break
			}
		}

		//if unsuccessful, handle the error
		if err != nil {
			handleErr(err, "MIDI device:"+string(id)+"cant receive a listener", true)
			if in.IsOpen() {
				//stop the device
				in.StopListening()
				in.Close()
			}
		} else {
			return err
		}
	}
	cliLog("MIDI Interf.", fmt.Sprintf("Listening to MIDI device %v", id), 0)
	return nil
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
// TO DO NOT WORKING YET
func initializeSavedDevices(devlist []helpers.InterfaceDevice, socket *websocket.Conn) (err error) {
	active := true
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
			listenErr := ListenMidi(device.HardwareID, device.HardwareID)

			if listenErr != nil {
				handleErr(listenErr, "unable to listen to device", true)
				active = false
			}

		} else {
			handleErr(nil, "unrecognized device type", true)
		}

		//update ui
		addUIdevice(device.DeviceType, device.HardwareName, device.FriendlyName, device.BindID, socket)

		//check if device exists
		// TO DO , disabled and enabled interfaces
		fmt.Println(active)

	}
	cliLog("Initialization", "Init completed. Devices loaded successfully", 0)
	return nil
}

func initializeUi(devlist []helpers.InterfaceDevice, socket *websocket.Conn) (err error) {
	//active := true
	for _, device := range devlist {
		//update ui
		if device.HardwareName != "" {
			//addUIdevice(device.DeviceType, device.HardwareName, device.FriendlyName, device.BindID, socket)
			addUIdevice(device.DeviceType, device.HardwareName, device.FriendlyName, device.BindID, socket)
		}

		//check if device exists
		// TO DO , disabled and enabled interfaces
	}
	cliLog("Initialization", "Init completed. Devices loaded successfully", 0)
	return nil
}
