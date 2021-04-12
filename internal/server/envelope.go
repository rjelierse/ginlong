package server

import (
	"bytes"
)

const (
	byteBegin = 0xA5
	byteEnd   = 0x15
)

const (
	addrBegin = iota
	addrLength
	_ // reserved
	_ // reserved
	addrMessageType
	addrProtocolVersion
	addrMessageCounter
	addrOriginator
	addrPayload = iota + 3
)

// Envelope wrap an incoming or outgoing message.
//
// It always has the following format:
// BEGIN LENGTH SIZE{2} TYPE VERSION COUNT SERIAL{4} ...payload... CHECKSUM END
//
// BEGIN:    Indicator a new message is starting.
// LENGTH:   The length of the payload field.
// SIZE:     Size for something? Always 4096.
// TYPE:     The message type.
// VERSION:  The protocol version? Depending on the response to a message it will switch measurements (and payload changes as well)
// COUNT:    The message counter for the current session.
// SERIAL:   The serial number of the data logging stick.
// PAYLOAD:  The message data.
// CHECKSUM: The checksum of all data after BEGIN (up to this point).
// END:      The terminating byte.
type Envelope []byte

// NewEnvelope constructs a new Envelope and validates its contents.
func NewEnvelope(data []byte) (Envelope, error) {
	envelope := Envelope(data)
	if err := envelope.Validate(); err != nil {
		return nil, err
	}
	return envelope, nil
}

// NewResponseEnvelope wraps a payload in Envelope.
//
// It sets the correct begin and end bytes and calculates the checksum for the payload.
func NewResponseEnvelope(payload []byte) Envelope {
	buf := bytes.NewBuffer(nil)

	buf.WriteByte(byteBegin)
	buf.Write(payload)
	buf.WriteByte(calculateChecksum(payload))
	buf.WriteByte(byteEnd)

	return buf.Bytes()
}

func (env Envelope) Length() uint8 {
	return env[addrLength]
}

func (env Envelope) MessageType() MessageType {
	return MessageType(env[addrMessageType])
}

func (env Envelope) ProtocolVersion() uint8 {
	return env[addrProtocolVersion]
}

func (env Envelope) MessageCounter() uint8 {
	return env[addrMessageCounter]
}

func (env Envelope) Originator() uint32 {
	return readUint32(env, addrOriginator)
}

func (env Envelope) Payload() []byte {
	return env[addrPayload : addrPayload+env.Length()]
}

func (env Envelope) Checksum() uint8 {
	return env[len(env)-2]
}

func (env Envelope) Validate() error {
	switch {
	case env[addrBegin] != byteBegin, env[len(env)-1] != byteEnd:
		return ErrInvalidMessage
	case env.Checksum() != calculateChecksum(env[1:len(env)-2]):
		return ErrChecksumMismatch
	}

	return nil
}

func calculateChecksum(message []byte) uint8 {
	var sum uint8
	for _, value := range message {
		sum = (sum + value) & 255
	}

	return sum
}
