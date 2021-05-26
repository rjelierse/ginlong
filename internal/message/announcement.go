package message

import (
	"net"

	"github.com/rjelierse/ginlong/internal/binary"
)

// Announcement contains the initial handshake by the data logger device.
type Announcement struct {
	HardwareAddress net.HardwareAddr
	IPAddress       net.IP
	FirmwareVersion string
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *Announcement) UnmarshalBinary(data []byte) error {
	const (
		addrFirmwareVersion = 0x13
		addrHardwareAddress = 0x3b
		addrNetworkAddress  = 0x41
	)

	m.HardwareAddress = data[addrHardwareAddress : addrHardwareAddress+6]
	m.IPAddress = net.ParseIP(binary.ReadString(data, addrNetworkAddress, 15))
	m.FirmwareVersion = binary.ReadString(data, addrFirmwareVersion, 15)
	return nil
}
