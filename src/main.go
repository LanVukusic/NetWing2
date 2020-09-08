package main

/* packr2 build -a -tags netgo -ldflags '-w -extldflags "-static"' -o myout2.exe main.go */

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

	"github.com/gobuffalo/packr/v2"
	"github.com/gomidi/connect"
	driver "github.com/gomidi/rtmididrv"
	"github.com/gorilla/websocket"
	"github.com/hypebeast/go-osc/osc"
)

// globals
var MIDIListenMode bool                                        // true:active interface, false:binding interface.
var listenDeviceType int                                       // type of listening.
var mappings map[helpers.InternalDevice]helpers.InternalOutput // array of mappings.
var OSClient osc.Client
var exec_pages []helpers.ExecWindow

//var OSClient2 osc.Client

// packeging
var boxView *packr.Box
var boxStatic *packr.Box

// main device array. supports max of 50 devices
var mainDeviceList []helpers.InterfaceDevice

// mapped interfaces and their bind ids.
var midi2idMappings []helpers.Bind2MIDI

// MIDI stuff
var drvMIDI connect.Driver

// Websocket and network
var upgrader websocket.Upgrader
var wsConnections []*websocket.Conn // supports max 30 clients

func main() {
	//init midi
	drvMIDI = nil
	MIDIListenMode = false
	mappings = make(map[helpers.InternalDevice]helpers.InternalOutput)

	exec_pages = []helpers.ExecWindow{}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// packages all static files in one binary
	boxView = packr.New("html_page", "./web/view")
	boxStatic = packr.New("static things", "./web/")

	// start services
	fmt.Println("locking Threads")
	runtime.LockOSThread()

	// start loopcheck to alert for disconnected devices
	go doEvery(2000*time.Millisecond, loopCheck)

	// create a http server to serve the UI both remote and to the local client
	fmt.Println("Starting webserver")
	runWebserver()
}

// helper functions and whatnot
func loopCheck() {
	devices, _ := getMIDIDevices(&drvMIDI)
	// check for inputs
	ins := devices.Ins
	//fmt.Println(devices, mainDeviceList)
	for i := range mainDeviceList {
		var enabled = false
		for _, allIterf := range ins {
			if mainDeviceList[i].HardwareName == allIterf.Name {
				//fmt.Println(mainDeviceList[i].HardwareName, allIterf.Name, "penis")
				enabled = true
				if !mainDeviceList[i].Active {
					cliLog("MIDI", fmt.Sprintf("Reconnecting device: %s", mainDeviceList[i].HardwareName), 0)
					ListenMidi(&drvMIDI, mainDeviceList[i].HardwareID, mainDeviceList[i].BindID, false)
					mainDeviceList[i].Active = true
				}
			}
		}
		if !enabled { // device was not found in active interfaces therefore it is disconnected
			if mainDeviceList[i].Active { // if device was not disconnected before this function call
				cliLog("MIDI", fmt.Sprintf("Device not connected: %s", mainDeviceList[i].HardwareName), 2)
				mainDeviceList[i].Active = false // mark it as disconnected device
			}
		}
	}
	return
}

func OSCstart(h string, oscIn *osc.Client) {
	//start OSC
	*oscIn = *osc.NewClient(h, 8000)
}

func doEvery(d time.Duration, f func()) {
	for range time.Tick(d) {
		f()
	}
}

func containsValue(m map[helpers.InternalDevice]helpers.InternalOutput, v helpers.InternalOutput) bool {
	for _, x := range m {
		if x == v {
			return true
		}
	}
	return false
}

func addUIdevice(deviceType int16, hName string, fName string, devID int, socket *websocket.Conn) {
	data := helpers.WSMsgTemplate{
		Event: "UiAddDevice",
		Data:  "{'ID':'" + strconv.Itoa(devID) + "','Hname':'" + hName + "','FriendlyName':'" + fName + "'}",
	}
	socket.WriteJSON(data)
}

func addUIFader(device interface{}, chn interface{}, MQChanel interface{}, socket *websocket.Conn, bcast bool) {
	// respond with permission to create UI fader
	data := helpers.MappingResponse{
		Event:     "MappingsResponse",
		Interface: 0, // 0= fader, 3 = exec
		DeviceID:  device.(int),
		ChannelID: chn.(byte),
		FaderID:   MQChanel.(int),
	}

	if bcast {
		broadcastMessage(data)
		return
	}
	socket.WriteJSON(data)
	return
}

func addUIExec(device interface{}, chn interface{}, execID interface{}, pageID interface{}, socket *websocket.Conn, bcast bool) {
	// respond with permission to create UI fader
	data := helpers.MappingResponse{
		Event:     "MappingsResponse",
		Interface: 3, // exec
		DeviceID:  device.(int),
		ChannelID: chn.(byte),
		ExecID:    execID.(int),
		ExecPage:  pageID.(int),
	}

	if bcast {
		broadcastMessage(data)
		return
	}
	socket.WriteJSON(data)
	return
}

func handleWSMessage(messageType int, p []byte, socket *websocket.Conn) {
	var raw map[string]interface{}
	err := json.Unmarshal(p, &raw)
	if err != nil {
		handleErr(err, fmt.Sprintf("Error while parsing JSON from %s", socket.LocalAddr().String()), true)
	}
	switch raw["event"] {
	// command for getting a list of active devices
	case "getMidiDevices":
		data, err := getMIDIDevices(&drvMIDI)
		if err != nil {
			handleErr(err, "Error while getting MIDI devices", true)
		} else {
			socket.WriteJSON(data)
		}
		break

	// command for setting midi mode to binding
	case "changeMIDImode": // changes midi mode to BIND
		fmt.Println("msg received BINDING active")
		MIDIListenMode = false
		listenDeviceType = int(raw["interface"].(float64))
		break

	case "addInterface":
		addInterface(int(raw["deviceType"].(float64)), fmt.Sprintf("%v", raw["HardwareName"]), fmt.Sprintf("%v", raw["FriendlyName"]), int(raw["inDevice"].(float64)))

	case "bindMIDIchannel":
		// future me will be thankful: https://blog.golang.org/maps

		//is channel used?
		tempKey := helpers.InternalDevice{
			InterfaceType: 0,
			DeviceID:      int(raw["device"].(float64)),
			ChannelID:     byte(raw["chn"].(float64)),
		}

		tempVal := helpers.InternalOutput{
			OutType: raw["extType"].(float64), // 3 is an EXEC, 0 is fader
			OutChan: int(raw["extChn"].(float64)),
			OutPage: int(raw["execPage"].(float64)),
			Fade:    raw["typeFader"].(bool),
		}

		//check the validity of mapping. two inputs cant be on same out and vice versa.
		for key, val := range mappings {
			if key.ChannelID == tempKey.ChannelID && key.DeviceID == tempKey.DeviceID {
				str := ""
				if val.OutType == 3 {
					str = "exec"
				}
				if val.OutType == 0 {
					str = "fader"
				}
				// this midi channel is in use
				cliLog("Mapping", "MIDI channel is already used on the "+str+": "+fmt.Sprintf("%v", val.OutChan)+", page:"+fmt.Sprintf("%v", val.OutPage), 1)
				return
			}
			if val.OutType == tempVal.OutType && val.OutChan == tempVal.OutChan && val.OutPage == tempVal.OutPage {
				cliLog("Mapping", "Can't map 2 inputs to same output interfaces", 1)
				return
			}
		}

		// create an internal mapping
		if int(raw["extType"].(float64)) == 0 { // FADER
			// fader
			mappings[tempKey] = tempVal
			addUIFader(int(raw["device"].(float64)), byte(raw["chn"].(float64)), int(raw["extChn"].(float64)), socket, true)
			break
		} else {
			if int(raw["extType"].(float64)) == 3 { // EXEC
				// add an entry to the mappings array
				mappings[tempKey] = tempVal
				addUIExec(int(raw["device"].(float64)), byte(raw["chn"].(float64)), int(raw["extChn"].(float64)), int(raw["execPage"].(float64)), socket, true)
				break
			}
		}

	case "addNewPage":
		temp := helpers.ExecWindow{
			Event:  "newExecPage",
			Page:   int(raw["page"].(float64)),
			Width:  int(raw["width"].(float64)),
			Height: int(raw["height"].(float64)),
		}
		exec_pages = append(exec_pages, temp)

		broadcastMessage(temp)

	case "restartOSC":
		//fmt.Println(raw["host"].(string))
		OSCstart(raw["host"].(string), &OSClient)
		cliLog("OSC", "Listening: "+raw["host"].(string), 1)

	case "removeMapping":
		removeMapping(int(raw["extChn"].(float64)), int(raw["execPage"].(float64)), raw["extType"].(float64))

	case "saveRequest":
		if raw["type"].(string) == "local" {
			temp := helpers.SingleData{
				Event: "saveReturn",
				JSN:   fmt.Sprintf(generateSaveJSON()),
			}
			socket.WriteJSON(temp)
			fmt.Println()
		}
	case "loadSave":
		mainDeviceList = []helpers.InterfaceDevice{}
		exec_pages = []helpers.ExecWindow{}
		midi2idMappings = []helpers.Bind2MIDI{}
		mappings = map[helpers.InternalDevice]helpers.InternalOutput{}

		for key, element := range raw["data"].(map[string]interface{}) {
			//fmt.Println(key, element)
			switch key {
			case "mappings": // all mappings map - needs to be proccesed seperately
				mappings = make(map[helpers.InternalDevice]helpers.InternalOutput)
				elem := element.([]interface{})
				for i := range elem {
					keyMaping := elem[i].(map[string]interface{})
					temp1 := keyMaping["Key"].(map[string]interface{})
					temp2 := keyMaping["Val"].(map[string]interface{})

					tempKey := helpers.InternalDevice{
						InterfaceType: int(temp1["InterfaceType"].(float64)),
						DeviceID:      int(temp1["DeviceID"].(float64)),
						ChannelID:     byte(temp1["ChannelID"].(float64)),
					}
					tempVal := helpers.InternalOutput{
						OutType: temp2["OutType"].(float64), // 3 is an EXEC, 0 is fader
						OutChan: int(temp2["OutChan"].(float64)),
						OutPage: int(temp2["OutPage"].(float64)),
						Fade:    temp2["Fade"].(bool),
					}
					mappings[tempKey] = tempVal
				}
				break

			case "exec_pages": // execs  list - jst copies all the values
				elem := element.([]interface{})
				for i := range elem {
					page := elem[i].(map[string]interface{})
					temp := helpers.ExecWindow{
						Event:  page["Event"].(string),
						Page:   int(page["Page"].(float64)),
						Width:  int(page["Width"].(float64)),
						Height: int(page["Height"].(float64)),
					}
					exec_pages = append(exec_pages, temp)
				}
				break

			case "mainDeviceList": // devices  list - reads all devices and tries to initialize them
				elem := element.([]interface{})
				for i := range elem {
					page := elem[i].(map[string]interface{})
					addInterface(int(page["DeviceType"].(float64)), page["HardwareName"].(string), page["FriendlyName"].(string), int(page["HardwareID"].(float64)))
				}
				break

			}
		}

		cliLog("Load", "Load successfull. Reload the page.", 0)
	}
}

func addInterface(devType int, hardwareNameI string, firendlyNameI string, inDev int) {
	if devType == 0 { // it is a MIDI device
		device := helpers.InterfaceDevice{
			Active:       true,
			BindID:       len(mainDeviceList), // that's the id for the device on the backend. not to be confused with midi ID
			DeviceType:   0,
			HardwareName: hardwareNameI,
			FriendlyName: firendlyNameI}

		mainDeviceList = append(mainDeviceList, device)

		// handle its inputs / outputs
		ListenMidi(&drvMIDI, inDev, len(mainDeviceList), true)
		sID := strconv.Itoa(len(mainDeviceList))

		//fmt.Println(sID)
		data := helpers.WSMsgTemplate{
			Event: "UiAddDevice",
			Data:  "{'ID':'" + sID + "','Hname':'" + hardwareNameI + "','FriendlyName':'" + firendlyNameI + "'}",
		}

		// update the UI
		broadcastMessage(data)
	}
}

func generateSaveJSON() string {
	// mappings
	jsonSave := "{ \"mappings\":["
	for key, element := range mappings {
		keyJSON, _ := json.Marshal(key)
		valJSON, _ := json.Marshal(element)
		jsonSave += fmt.Sprintf("{\"Key\": %s, \"Val\": %s},", keyJSON, valJSON)
	}
	jsonSave = strings.TrimRight(jsonSave, ",")

	// exec_pages
	jsonSave += "], \"exec_pages\":"
	execPages, _ := json.Marshal(exec_pages)
	jsonSave += fmt.Sprintf("%s", execPages)

	// mainDeviceList
	jsonSave += ", \"mainDeviceList\":"
	mainDeviceJSON, _ := json.Marshal(mainDeviceList)
	jsonSave += fmt.Sprintf("%s", mainDeviceJSON)

	jsonSave += "}"
	return jsonSave
}

func removeMapping(channel int, page int, typ float64) {
	for key, element := range mappings {
		if element.OutType == typ && element.OutPage == page && element.OutChan == channel {
			delete(mappings, key)
			removeMsg := helpers.MappingRemove{
				Event:   "removeMappingRes",
				Page:    page,
				Channel: channel,
				Type:    typ,
			}
			broadcastMessage(removeMsg)
			return
		}
	}

	cliLog("mapping", fmt.Sprintf("Can't remove mapping MIDI:%d Page:%d", channel, page), 1)
	return
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
				if strings.Contains(err.Error(), "close 1001 (going away)") {
					cliLog("WebServer", "Disconnect! Bye bye "+wsConnection.LocalAddr().String(), 1)
					for i := 0; i < len(wsConnections); i++ {
						if wsConnections[i] == wsConnection {
							// we have found the faulty connection, so we remove it to prevent faulty broadcasting
							wsConnections = append((wsConnections)[:i], wsConnections[i+1:]...)
							return
						}
					}
				}
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
		fmt.Fprintf(w, "Testing. go to <a href=\"/ui\"> this thing</a>")
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

func getMIDIDevices(drv *connect.Driver) (outDevices helpers.MidiPackage, err error) {
	outDevices.Event = "refreshMidiRet"

	// does not get inputs, aka an error
	if *drv == nil { // if driver does not exist yet, we create the new driver
		cliLog("MIDI", "No active Midi driver", 1)
		*drv, err = driver.New()
		if err != nil {
			handleErr(err, "Failed creation of MIDI device driver.", true)
			return outDevices, err
		}
		cliLog("MIDI", "MIDI driver activated", 0)
	}

	var tempDrv = *drv // use the drivers value and not the pointer from that point on

	// gets the inputs
	ins, err := tempDrv.Ins()
	if err != nil { // handle the error
		handleErr(err, "Error receiving midi device list", true)
		return outDevices, err
	}

	for _, el := range ins {
		outDevices.Ins = append(outDevices.Ins, helpers.MidiDevice{
			// it splits the string by spaces, removeslast slice (the number) and joins it back together to form a string
			NameWithID: el.String(),
			Name:       strings.Join(strings.Fields(el.String())[:len(strings.Fields(el.String()))-1], ""),
			ID:         el.Number()})
	}

	outs, err := tempDrv.Outs()
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

	MIDItype := int(in[0])
	if MIDItype == 192 {
		cliLog("MIDI", "Cant map to 'CONTROLL CHANGE' messages", 2)
		return
	}
	MIDIchannel := in[1]
	MIDIvalue := in[2]

	if MIDIListenMode {
		// input has an active binding
		tempIn := helpers.InternalDevice{
			InterfaceType: 0, // MIDI = 0,
			DeviceID:      deviceID,
			ChannelID:     MIDIchannel,
		}
		tempOut, exists := mappings[tempIn]

		if exists {

			switch tempOut.OutType {
			case 0:

				temp := helpers.FaderUpdate{
					Event:   "UpdateFader",
					Type:    0,
					FaderID: tempOut.OutChan,
					Value:   MIDIvalue,
				}
				broadcastMessage(temp)
				msg := osc.NewMessage("/pb/" + fmt.Sprintf("%d", tempOut.OutChan))
				msg.Append(fmt.Sprintf("%d", int((int(MIDIvalue) * 100 / 127))))
				OSClient.Send(msg)
				break
			case 3:
				// check if it's a button

				var out int
				if !tempOut.Fade {
					//check control type
					switch MIDItype {
					case 144:
						// note on - we turn button on
						out = 127
						break
					case 128:
						// note off -we turn button off
						out = 0
						break
					case 176:
						// if value is grater than 0, button is on
						if MIDIvalue != 0 {
							out = 127
							break
						}
					}
				} else {
					out = int(MIDIvalue)
				}

				temp := helpers.ExecUpdate{
					Event:    "UpdateFader",
					Type:     3,
					FaderID:  tempOut.OutChan,
					PageID:   tempOut.OutPage,
					FadeType: tempOut.Fade,
					Value:    out,
				}

				broadcastMessage(temp)
				msg := osc.NewMessage("/exec/" + fmt.Sprintf("%d", tempOut.OutPage) + "/" + fmt.Sprintf("%d", tempOut.OutChan))
				msg.Append(fmt.Sprintf("%d", int((out * 100 / 127))))
				//msg.Append(int((out * 100 / 127)))
				OSClient.Send(msg)
				break
			}
		} else {
			// input has no active binding
			cliLog("MIDI", fmt.Sprintf("Type: %v, Chn: %v, Val: %v, Device: %v", MIDItype, MIDIchannel, MIDIvalue, int(deviceID)), 0)
		}
	} else {
		// binding interface mode
		temp := helpers.MIDILearnMessage{
			Event:     "learnMidiRet",
			Interf:    listenDeviceType,
			DeviceID:  deviceID,
			ChannelID: MIDIchannel,
		}
		broadcastMessage(temp)
		MIDIListenMode = true // exit midi mapping mode
	}

}

func cliLog(cause string, body string, threatLevel int) {
	alert := helpers.CliMsg{
		Event:       "cli",
		Cause:       cause,
		Body:        body,
		ThreatLevel: threatLevel,
	}
	broadcastMessage(alert)
	fmt.Println(alert)
	return
}

func handleErr(err error, msg string, bcast bool) {
	fmt.Println(err, msg)
	if bcast {
		cliLog(err.Error(), msg, 2)
	}

}

//ListenMidi checks availability and ATTACHES a MIDI listener to the device
func ListenMidi(drv *connect.Driver, id int, bind int, newDevice bool) (err error) {
	// first check for existence of MIDI driver
	if *drv == nil { // if driver does not exist yet, we create the new driver
		cliLog("MIDI", "No active Midi driver  ... activating", 1)
		*drv, err = driver.New()
		if err != nil {
			handleErr(err, "Failed creation of MIDI device driver.", true)
			return err
		}
		cliLog("MIDI", "MIDI driver activated", 0)
	}

	// this line gets the device "id" from the driver, opens it and returns the active device
	in, err := connect.OpenIn(drvMIDI, id, "")

	//handles the potential error
	if err != nil { // if error exists, handle it and return
		handleErr(err, "MIDI device:"+string(id)+"is unavailable", true)

	}

	// add device to monitroring array
	device := helpers.Bind2MIDI{ // assign it the object with the device
		BindID:    bind,
		MidiPort:  in,
		WasOnline: true}

	midi2idMappings = append(midi2idMappings, device)

	//if the device is successfully opened it tries to attach a listener
	if err == nil {

		// set a listener to the device
		err := in.SetListener(handleMidiEvent)

		// check for potential errors
		if err != nil { // if unsuccessful, handle the error
			handleErr(err, "MIDI device:"+string(id)+"cant receive a listener", true)
			if in.IsOpen() {
				//stop the device
				in.StopListening()
				in.Close()
			}
			return err
		}
		cliLog("MIDI", fmt.Sprintf("Listening to MIDI device %v", id), 0)
		return nil

	}
	cliLog("MIDI", "Midi device is not open", 2)
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

func initializeUi(devlist []helpers.InterfaceDevice, socket *websocket.Conn) (err error) {
	//active := true
	for _, device := range devlist {
		//update ui
		if device.HardwareName != "" {
			//addUIdevice(device.DeviceType, device.HardwareName, device.FriendlyName, device.BindID, socket)
			addUIdevice(device.DeviceType, device.HardwareName, device.FriendlyName, device.BindID, socket)
		}

		//check if device exists
	}

	// update pages
	for _, s := range exec_pages {
		fmt.Println(s)
		temp := helpers.ExecWindow{
			Event:  "newExecPage",
			Page:   s.Page,
			Width:  s.Width,
			Height: s.Height,
		}
		socket.WriteJSON(temp)
	}

	for k, v := range mappings {
		//fmt.Printf("key[%s] value[%s]\n", k, v)
		switch v.OutType {
		case 0:
			addUIFader(k.DeviceID, k.ChannelID, v.OutChan, socket, false)
			break
		case 3:
			addUIExec(k.DeviceID, k.ChannelID, v.OutChan, v.OutPage, socket, false)
			break
		}

	}

	cliLog("Initialization", "Init completed. Devices loaded successfully", 0)
	return nil
}
