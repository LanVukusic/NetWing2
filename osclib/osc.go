package osclib

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"
)

// OscClient is a client that serves OSC data
var OscClient osc.Client

// SendOSC is a debug function
func SendOSC(channel int, value int) {
	msg := osc.NewMessage("/")
	msg.Append("Chan: ")
	msg.Append(int32(channel))
	msg.Append("Value: ")
	msg.Append(int32(value))
	OscClient.Send(msg)
	fmt.Println(channel, value)
}

// StartOSCserver sdasd
func StartOSCserver() {
	OscClient = *osc.NewClient("192.168.1.8", 1234)
}
