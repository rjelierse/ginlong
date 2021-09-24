// Package messages contains the types for messages passing through the system.
package messages

import (
	"encoding/binary"
	"math"
)

// MessageType indicates the type of message being sent.
type MessageType uint8

// MessageTypes from monitoring to logging stick.
const (
	ResponseTypeAnnounce MessageType = 0x11 + iota
	ResponseTypeMeasurement
	_ // reserved
	_ // reserved
	_ // reserved
	_ // reserved
	ResponseTypePing
	ResponseType48
)

// MessageTypes from logging stick to monitoring.
const (
	MessageTypeAnnounce MessageType = 0x41 + iota
	MessageTypeMeasurement
	_ // reserved
	_ // reserved
	_ // reserved
	_ // reserved
	MessageTypePing
	MessageType48
)

// Envelope contains the information on a message sent between the logging stick and the monitor system.
type Envelope struct {
	MessageType       MessageType
	MonitoringCounter uint8
	LoggerCounter     uint8
	LoggerSerial      uint32
	Checksum          uint8
	Message           Message
}

// NewEnvelope constructs an Envelope from data.
func NewEnvelope(data []byte) Envelope {
	return Envelope{
		MessageType:       MessageType(data[0x04]),
		MonitoringCounter: data[0x05],
		LoggerCounter:     data[0x06],
		LoggerSerial:      binary.LittleEndian.Uint32(data[0x07:]),
		Checksum:          data[len(data)-2],
		Message:           data[0x0b : 0x0b+data[0x01]],
	}
}

// Message wraps the message data contained in the envelope.
type Message []byte

// Read16bit decodes an unsigned 16-bit value from the message and returns it with the given scale.
func (msg Message) Read16bit(offset int, scale float64) float64 {
	if len(msg) < offset+1 {
		return math.Inf(1)
	}
	v := binary.LittleEndian.Uint16(msg[offset:])
	return float64(v) * scale
}

// Read32bit decodes an unsigned 32-bit value from the message and returns it with the given scale.
func (msg Message) Read32bit(offset int, scale float64) float64 {
	if len(msg) < offset+3 {
		return math.Inf(1)
	}
	v := binary.LittleEndian.Uint32(msg[offset:])
	return float64(v) * scale
}

// ReadString decodes a string value from the message.
func (msg Message) ReadString(offset int, length int) string {
	return string(msg[offset : offset+length])
}
