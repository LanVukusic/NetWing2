package main

import (
	"fmt"

	"github.com/gomidi/connect"
	driver "github.com/gomidi/rtmididrv"
)

func main() {
	// Create a new midi driver and close it with "defer" when the program terminates
	var drv connect.Driver
	drv, _ = driver.New()

	//get the input devices
	fmt.Println(getMIDIDevices(drv))

	// open the device in index "0", without using device id
	in, _ := connect.OpenIn(drv, 0, "")

	// open the device by passing the name, and giving an id >= 0
	//example connect.OpenIn(drv, 1, "AKAI 0")

	//in, _ := connect.OpenIn(drv, 1, getMIDIDevices(drv)[0].String())

	in.SetListener(handleMIDIevent)

	// we can keep the thread running by keeping the while loop
	for {

	}
}

// gets all the midi devices from the driver and returns them
func getMIDIDevices(drv connect.Driver) (outs []connect.In) {
	//gets the inputs
	ins, _ := drv.Ins()

	/* // that is the way to get the outputs
	ins, _ := drv.Outs() */

	for _, el := range ins {
		outs = append(outs, el)
	}
	return outs
}

func handleMIDIevent(data []byte, deltaMicroseconds int64) {
	fmt.Printf("Midi message on: %v value: %v \n", data[1], data[2])
}
