package messages

import (
	"encoding/binary"
	"fmt"
	"time"
)

type (
	// Measurement contains the data for a single measurement point on an inverter.
	Measurement struct {
		Timestamp time.Time
		Inverter  Inverter
		PV1       SolarArray
		PV2       SolarArray
		Output    Output
		Grid      Grid
	}

	// SolarArray represents a solar array input on the inverter.
	SolarArray struct {
		Voltage float64
		Current float64
	}

	// Output contains the statistics for the output generated by the inverter.
	Output struct {
		Today     float64
		Yesterday float64
		ThisMonth float64
		LastMonth float64
		ThisYear  float64
		Lifetime  float64
	}

	// Grid represents the status of the electrical grid to which the inverter is connected.
	Grid struct {
		CurrentP1 float64
		CurrentP2 float64
		CurrentP3 float64
		Voltage   float64
		Frequency float64

		ApparentPower float64
		PowerFactor   float64
	}

	// Inverter contains the information on the inverter that recorded the measurement.
	Inverter struct {
		Serial      string
		Sensor      string
		Model       string
		DSPVersion  string
		HCIVersion  string
		Temperature float64
	}
)

const (
	offsetSensor             = 0x01 + (iota * 2) // 0x01 = 1
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	offsetSerialNumber                           // 0x15 = 21
	_                                            // reserved
	_                                            // reserved
	_                                            // reserved
	_                                            // reserved
	_                                            // reserved
	_                                            // reserved
	_                                            // reserved
	offsetInverterTemp                           // 0x25 = 37
	offsetArray1Voltage                          // 0x27
	offsetArray2Voltage                          // 0x29
	offsetArray1Current                          // 0x2b
	offsetArray2Current                          // 0x2d
	_                                            // unknown
	_                                            // unknown
	offsetGridPhase1Current                      // 0x33 = 51
	offsetGridPhase2Current                      // 0x35
	offsetGridPhase3Current                      // 0x37
	offsetGridVoltage                            // 0x39
	offsetGridFrequency                          // 0x3b
	_                                            // unknown
	_                                            // unknown
	offsetOutputToday                            // 0x41 = 65
	_                                            // reserved - previous offset is 32bit
	offsetOutputTotal                            // 0x45 = 69
	_                                            // reserved - previous offset is 32bit
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	offsetVBus                                   // 0x61 = 97
	offsetHVBus                                  // 0x63
	offsetPowerLimit                             // 0x65
	offsetGridPowerFactor                        // 0x67
	offsetArrayTotalPower                        // 0x69
	_                                            // unknown
	offsetOutputThisMonth                        // 0x6d = 109
	_                                            // reserved - previous offset is 32bit
	offsetOutputLastMonth                        // 0x71
	_                                            // reserved - previous offset is 32bit
	offsetOutputYesterday                        // 0x75
	offsetOutputThisYear                         // 0x77
	_                                            // reserved - previous offset is 32bit
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	_                                            // unknown
	offsetGridApparentPower                      // 0x83 = 131
	_                                            // unknown
	offsetTimestampYear                          // 0x87 = 135
	offsetTimestampMonth                         // 0x89
	offsetTimestampDay                           // 0x8b
	offsetTimestampHour                          // 0x8d
	offsetTimestampMinute                        // 0x8f
	offsetTimestampSecond                        // 0x91
	offsetInverterModel                          // 0x93
	offsetInverterDSPVersion                     // 0x95
	offsetInverterHCIVersion                     // 0x97
)

const (
	scaleNone      = 1.0
	scaleTenths    = 0.1
	scaleHundreds  = 0.01
	scaleThousands = 0.001
)

func NewMeasurement(msg Message) (Measurement, error) {
	ts, err := time.Parse("06-01-02 15:04:05 -0700", fmt.Sprintf(
		"%02d-%02d-%02d %02d:%02d:%02d +0200",
		binary.LittleEndian.Uint16(msg[offsetTimestampYear:]),
		binary.LittleEndian.Uint16(msg[offsetTimestampMonth:]),
		binary.LittleEndian.Uint16(msg[offsetTimestampDay:]),
		binary.LittleEndian.Uint16(msg[offsetTimestampHour:]),
		binary.LittleEndian.Uint16(msg[offsetTimestampMinute:]),
		binary.LittleEndian.Uint16(msg[offsetTimestampSecond:]),
	))
	if err != nil {
		return Measurement{}, err
	}

	return Measurement{
		Timestamp: ts,
		Inverter: Inverter{
			Sensor:      fmt.Sprintf("%04x", binary.LittleEndian.Uint16(msg[offsetSensor:])),
			Serial:      msg.ReadString(offsetSerialNumber, 16),
			Model:       fmt.Sprintf("%02x", msg[offsetInverterModel]),
			DSPVersion:  fmt.Sprintf("%04x", binary.LittleEndian.Uint16(msg[offsetInverterDSPVersion:])),
			HCIVersion:  fmt.Sprintf("%04x", binary.LittleEndian.Uint16(msg[offsetInverterHCIVersion:])),
			Temperature: msg.Read16bit(offsetInverterTemp, scaleTenths),
		},
		PV1: SolarArray{
			Voltage: msg.Read16bit(offsetArray1Voltage, scaleTenths),
			Current: msg.Read16bit(offsetArray1Current, scaleTenths),
		},
		PV2: SolarArray{
			Voltage: msg.Read16bit(offsetArray2Voltage, scaleTenths),
			Current: msg.Read16bit(offsetArray2Current, scaleTenths),
		},
		Output: Output{
			Today:     msg.Read32bit(offsetOutputToday, scaleHundreds),
			Yesterday: msg.Read16bit(offsetOutputYesterday, scaleTenths),
			ThisMonth: msg.Read32bit(offsetOutputThisMonth, scaleNone),
			LastMonth: msg.Read32bit(offsetOutputLastMonth, scaleNone),
			ThisYear:  msg.Read32bit(offsetOutputThisYear, scaleNone),
			Lifetime:  msg.Read32bit(offsetOutputTotal, scaleTenths),
		},
		Grid: Grid{
			CurrentP1:     msg.Read16bit(offsetGridPhase1Current, scaleTenths),
			CurrentP2:     msg.Read16bit(offsetGridPhase2Current, scaleTenths),
			CurrentP3:     msg.Read16bit(offsetGridPhase3Current, scaleTenths),
			Voltage:       msg.Read16bit(offsetGridVoltage, scaleTenths),
			Frequency:     msg.Read16bit(offsetGridFrequency, scaleHundreds),
			ApparentPower: msg.Read16bit(offsetGridApparentPower, scaleNone),
			PowerFactor:   msg.Read16bit(offsetGridPowerFactor, scaleThousands),
		},
	}, nil
}
