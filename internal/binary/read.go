package binary

import (
	"encoding/binary"
	"strings"
)

// ReadFloat64 unpacks a uint16 value from data and applies scale to return a float64.
func ReadFloat64(data []byte, start, scale int) float64 {
	v := ReadUint16(data, start)

	return float64(v) / float64(scale)
}

// ReadUint16 unpacks a uint16 value from data.
func ReadUint16(data []byte, start int) uint16 {
	return binary.LittleEndian.Uint16(data[start : start+2])
}

// ReadUint32 unpacks a uint32 value from data.
func ReadUint32(data []byte, start int) uint32 {
	return binary.LittleEndian.Uint32(data[start : start+4])
}

// ReadString unpacks a string value from data.
func ReadString(data []byte, start, length int) string {
	return strings.Trim(string(data[start:start+length]), "\x00")
}
