package flightpanels

import (
	"log"
	"github.com/google/gousb"
)

type Switch uint

type SwitchState struct {
	Switch Switch
	Value  uint
}

type PanelReader interface {
	noZeroSwitch(i Switch) bool
}


func readSwitches(panel PanelReader, inEndpoint *gousb.InEndpoint, c chan SwitchState) {
	var data [3]byte
	var state uint64
	var newState uint64

	stream, err := inEndpoint.NewStream(3, 1)
	if err != nil {
		log.Fatalf("Could not create read stream: %v", err)
	}
	defer stream.Close()

	for {
		_, err := stream.Read(data[:])
		if err != nil {
			log.Fatalf("Read error: %v", err)
		}
		newState = uint64(data[0]) | uint64(data[1])<<8 | uint64(data[2])<<16
		changed := state ^ newState
		state = newState
		for i := Switch(0); i < 24 ; i++ {
			if (changed>>i)&1 == 1 {
				val := uint(state >> i & 1)
				if val == 0 && panel.noZeroSwitch(i) {
					continue
				}
				c <- SwitchState{i, val}
			}
		}
	}
}