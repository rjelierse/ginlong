package message

// Ping is sent as a keep-alive.
type Ping struct{}

// UnmarshalBinary implements encoding.BinaryUnmarshaler.
func (m *Ping) UnmarshalBinary(data []byte) error {
	return nil
}
