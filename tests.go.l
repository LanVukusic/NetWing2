package main

import (
	"fmt"
	"reflect"

	driver "github.com/gomidi/rtmididrv"
)

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}

// This example expects the first input and output port to be connected
// somehow (are either virtual MIDI through ports or physically connected).
// We write to the out port and listen to the in port.
func main() {
	drv, err := driver.New()
	must(err)
	fmt.Print(reflect.TypeOf(drv))

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	fmt.Print(ins)
	/* for _, el := range ins {
		fmt.Println(reflect.TypeOf(el))
		fmt.Print(el.String(), el.Number())
	} */

}

/*
func getMIDIdevices(drv driver) (device rtmididrv.in) {

}
*/
