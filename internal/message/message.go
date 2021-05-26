package message

import (
	"fmt"
)

// Type indicates the message type.
type Type uint8

// The different message types - see PROTOCOL.md for more information.
// TypeAnnouncement is sent when the data logger connects.
// TypeMeasurement is sent when the data logger exports a measurement set for an inverter.
// TypeAccessPointInfo is an informative message containing the network settings for the data logger.
// TypePing is a keep-alive message.
// TypeUnk2 is an unknown message.
const (
	TypeAnnouncement Type = 0x41 + iota
	TypeMeasurement
	TypeAccessPointInfo
	_
	_
	_
	TypePing
	TypeUnk2
)

// String implements fmt.Stringer.
func (mt Type) String() string {
	switch mt {
	case TypeAnnouncement:
		return "announce"
	case TypeMeasurement:
		return "measurement"
	case TypePing:
		return "ping"
	case TypeAccessPointInfo:
		return "accesspoint"
	default:
		return fmt.Sprintf("unknown#%x", uint8(mt))
	}
}

// ResponseType returns the expected response type for this type of message.
func (mt Type) ResponseType() uint8 {
	return uint8(mt) - 0x30
}
