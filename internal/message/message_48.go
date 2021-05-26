package message

import "github.com/rjelierse/ginlong/internal/binary"

// Message48 (unknown)
type Message48 struct {
	FirmwareVersion string
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *Message48) UnmarshalBinary(data []byte) error {
	const (
		addrFirmwareVersion = 0x14
	)
	m.FirmwareVersion = binary.ReadString(data, addrFirmwareVersion, 15)
	return nil
}
