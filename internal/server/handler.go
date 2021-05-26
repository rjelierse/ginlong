package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/rjelierse/ginlong/internal/message"

	"github.com/rs/zerolog"
)

const (
	ProtocolV2 = 0x02 + iota
	ProtocolV3
)

type handler struct {
	logger          zerolog.Logger
	defaultProtocol uint8
}

func (h *handler) handleMessage(ctx context.Context, envelope Envelope) (Envelope, error) {
	switch envelope.MessageType() {
	case message.TypeAnnouncement:
		return h.handleAnnounce(ctx, envelope)
	case message.TypeMeasurement:
		return h.handleMeasurement(ctx, envelope)
	case message.TypePing:
		return h.handlePing(ctx, envelope)
	case message.TypeAccessPointInfo:
		return h.handleAccessPointInfo(ctx, envelope)
	case message.TypeUnk2:
		return h.handleMessage48(ctx, envelope)
	default:
		return newResponse(envelope, h.defaultProtocol, newResponsePayload(0, 0)), fmt.Errorf("%w: %s", ErrUnsupportedMessage, envelope.MessageType())
	}
}

func (h *handler) handleAnnounce(ctx context.Context, envelope Envelope) (Envelope, error) {
	var msg message.Announcement
	if err := msg.UnmarshalBinary(envelope.Payload()); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}
	return newResponse(envelope, h.defaultProtocol, newResponsePayload(0x02, 0x01)), nil
}

func (h *handler) handleMeasurement(ctx context.Context, envelope Envelope) (Envelope, error) {
	var msg message.Measurement
	if err := msg.UnmarshalBinary(envelope.Payload()); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}

	h.logger.Debug().
		Interface("measurement", msg).
		Interface("fields", msg.Fields).
		Msg("Measurement received")

	return newResponse(envelope, h.defaultProtocol, newResponsePayload(0x01, 0x01)), nil
}

func (h *handler) handlePing(ctx context.Context, envelope Envelope) (Envelope, error) {
	var msg message.Ping
	if err := msg.UnmarshalBinary(envelope.Payload()); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}
	return newResponse(envelope, h.defaultProtocol, newResponsePayload(0x00, 0x00)), nil
}

func (h *handler) handleAccessPointInfo(ctx context.Context, envelope Envelope) (Envelope, error) {
	var msg message.AccessPointInfo
	if err := msg.UnmarshalBinary(envelope.Payload()); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}
	return newResponse(envelope, h.defaultProtocol, newResponsePayload(0x01, 0x01)), nil
}

func (h *handler) handleMessage48(ctx context.Context, envelope Envelope) (Envelope, error) {
	var msg message.Message48
	if err := msg.UnmarshalBinary(envelope.Payload()); err != nil {
		return nil, fmt.Errorf("failed to parse payload: %w", err)
	}
	return newResponse(envelope, h.defaultProtocol, newResponsePayload(0x01, 0x01)), nil
}

func newResponsePayload(usr1, usr2 byte) []byte {
	buf := bytes.NewBuffer(nil)
	buf.WriteByte(usr1)
	buf.WriteByte(usr2)
	_ = binary.Write(buf, binary.LittleEndian, uint32(time.Now().Unix()))
	buf.WriteByte(0x78)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	buf.WriteByte(0x00)
	return buf.Bytes()
}

func newResponse(envelope Envelope, protocolVersion uint8, payload []byte) Envelope {
	responseType := uint8(envelope.MessageType()) - 0x30

	buf := bytes.NewBuffer(nil)
	buf.WriteByte(uint8(len(payload)))
	buf.WriteByte(0x00)
	buf.WriteByte(0x10)
	buf.WriteByte(responseType)
	buf.WriteByte(protocolVersion)
	buf.WriteByte(envelope.MessageCounter())
	_ = binary.Write(buf, binary.LittleEndian, envelope.Originator())
	buf.Write(payload)

	return NewResponseEnvelope(buf.Bytes())
}
