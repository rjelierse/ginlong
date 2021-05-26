package message

import "github.com/rjelierse/ginlong/internal/binary"

// AccessPointInfo contains the data logger wifi settings.
type AccessPointInfo struct {
	SSID string
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *AccessPointInfo) UnmarshalBinary(data []byte) error {
	const (
		addrSSID = 0x0f
	)
	m.SSID = binary.ReadString(data, addrSSID, 30)
	return nil
}
