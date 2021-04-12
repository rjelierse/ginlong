package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

type MessageType uint8

func (mt MessageType) String() string {
	switch mt {
	case MessageTypeAnnounce:
		return "announce"
	case MessageTypeMeasurement:
		return "measurement"
	case MessageTypePing:
		return "ping"
	case MessageTypeAccessPointInfo:
		return "accesspoint"
	default:
		return fmt.Sprintf("unknown#%x", uint8(mt))
	}
}

func (mt MessageType) ResponseType() uint8 {
	return uint8(mt) - 0x30
}

const (
	MessageTypeAnnounce MessageType = 0x41 + iota
	MessageTypeMeasurement
	MessageTypeAccessPointInfo
	_
	_
	_
	MessageTypePing
	MessageTypeUnk2
)

// Announcement contains the initial handshake by the data logger device.
type Announcement struct {
	MacAddress      net.HardwareAddr
	IpAddress       net.IP
	FirmwareVersion string
}

func (m *Announcement) UnmarshalBinary(data []byte) error {
	const (
		addrFirmwareVersion = 0x13
		addrHardwareAddress = 0x3b
		addrNetworkAddress  = 0x41
	)

	m.MacAddress = data[addrHardwareAddress : addrHardwareAddress+6]
	m.IpAddress = net.ParseIP(readString(data, addrNetworkAddress, 15))
	m.FirmwareVersion = readString(data, addrFirmwareVersion, 15)
	return nil
}

// Measurement contains the data as measured on the inverter.
type Measurement struct {
	Inverter struct {
		SerialNumber  string
		Model         uint16
		FirmwareMain  uint16
		FirmwareSlave uint16
	}
}

func (m *Measurement) UnmarshalBinary(data []byte) error {
	const (
		addrInverterSerial  = 0x15
		addrInverterModel   = 0x93
		addrInverterFWSlave = 0x95
		addrInverterFWMain  = 0x97
	)
	m.Inverter.SerialNumber = readString(data, addrInverterSerial, 15)
	m.Inverter.Model = readUint16(data, addrInverterModel)
	m.Inverter.FirmwareMain = readUint16(data, addrInverterFWMain)
	m.Inverter.FirmwareSlave = readUint16(data, addrInverterFWSlave)
	return nil
}

// AccessPointInfo contains the data logger wifi settings.
type AccessPointInfo struct {
	SSID string
}

func (m *AccessPointInfo) UnmarshalBinary(data []byte) error {
	const (
		addrSSID = 0x0f
	)
	m.SSID = readString(data, addrSSID, 30)
	return nil
}

// Ping is sent as a keep-alive.
type Ping struct{}

func (m *Ping) UnmarshalBinary(data []byte) error {
	return nil
}

// Message48 (unknown)
type Message48 struct {
	FirmwareVersion string
}

func (m *Message48) UnmarshalBinary(data []byte) error {
	const (
		addrFirmwareVersion = 0x14
	)
	m.FirmwareVersion = readString(data, addrFirmwareVersion, 15)
	return nil
}

func readUint16(data []byte, start int) uint16 {
	return binary.LittleEndian.Uint16(data[start : start+2])
}

func readUint32(data []byte, start int) uint32 {
	return binary.LittleEndian.Uint32(data[start : start+4])
}

func readString(data []byte, start, length int) string {
	return strings.Trim(string(data[start:start+length]), "\x00")
}
