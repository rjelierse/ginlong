package message

import (
	"fmt"

	"github.com/rjelierse/ginlong/internal/binary"
)

// Measurement contains the data as measured on the inverter.
type Measurement struct {
	Inverter struct {
		SerialNumber  string
		Model         uint16
		FirmwareMain  uint16
		FirmwareSlave uint16
	}

	Temperature float64

	Array1 struct {
		Voltage float64
		Current float64
	}

	Array2 struct {
		Voltage float64
		Current float64
	}

	Grid struct {
		Phase1    float64
		Phase2    float64
		Phase3    float64
		Voltage   float64
		Frequency float64
	}

	Fields map[string]uint16
}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *Measurement) UnmarshalBinary(data []byte) error {
	const (
		addrInverterSerial  = 0x15
		addrInverterModel   = 0x93
		addrInverterFWSlave = 0x95
		addrInverterFWMain  = 0x97
		addrTemperature     = 0x25
		addrVpv1            = 0x27
		addrIpv1            = 0x2b
		addrVpv2            = 0x29
		addrIpv2            = 0x2d
		addrIac1            = 0x33
		addrIac2            = 0x35
		addrIac3            = 0x37
		addrVac             = 0x39
		addrFac             = 0x3b
	)
	m.Inverter.SerialNumber = binary.ReadString(data, addrInverterSerial, 15)
	m.Inverter.Model = binary.ReadUint16(data, addrInverterModel)
	m.Inverter.FirmwareMain = binary.ReadUint16(data, addrInverterFWMain)
	m.Inverter.FirmwareSlave = binary.ReadUint16(data, addrInverterFWSlave)

	m.Temperature = binary.ReadFloat64(data, addrTemperature, 10)

	m.Array1.Voltage = binary.ReadFloat64(data, addrVpv1, 10)
	m.Array1.Current = binary.ReadFloat64(data, addrIpv1, 10)

	m.Array2.Voltage = binary.ReadFloat64(data, addrVpv2, 10)
	m.Array2.Current = binary.ReadFloat64(data, addrIpv2, 10)

	m.Grid.Phase1 = binary.ReadFloat64(data, addrIac1, 10)
	m.Grid.Phase2 = binary.ReadFloat64(data, addrIac2, 10)
	m.Grid.Phase3 = binary.ReadFloat64(data, addrIac3, 10)
	m.Grid.Voltage = binary.ReadFloat64(data, addrVac, 10)
	m.Grid.Frequency = binary.ReadFloat64(data, addrFac, 100)

	m.Fields = make(map[string]uint16)
	for addr := addrTemperature; addr < addrInverterModel; addr += 2 {
		m.Fields[fmt.Sprintf("0x%02x", addr)] = binary.ReadUint16(data, addr)
	}

	return nil
}
