package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/rs/zerolog"
)

const (
	bufSize = 4096
)

var (
	ErrChecksumMismatch   = errors.New("message checksum does not match data")
	ErrInvalidMessage     = errors.New("invalid message")
	ErrUnsupportedMessage = errors.New("unsupported message type")
)

type Server struct {
	logger     zerolog.Logger
	handler    *handler
	wg         sync.WaitGroup
	shutdownFn context.CancelFunc
}

func New(logger zerolog.Logger) *Server {
	return &Server{
		logger: logger,
		handler: &handler{
			logger:          logger,
			defaultProtocol: ProtocolV3,
		},
	}
}

func (s *Server) Listen(address string) error {
	lc := net.ListenConfig{}

	var ctx context.Context
	ctx, s.shutdownFn = context.WithCancel(context.Background())
	socket, err := lc.Listen(ctx, "tcp4", address)
	if err != nil {
		return err
	}

	s.logger.Info().Str("address", socket.Addr().String()).Msg("Listener started, waiting for connections")

	go s.startAcceptingConnections(socket)

	return nil
}

func (s *Server) Shutdown() {
	s.logger.Info().Msg("Shutting down")
	s.shutdownFn()

	// Wait until all connections are drained
	s.wg.Wait()
}

func (s *Server) startAcceptingConnections(listener net.Listener) {
	for {
		socket, err := listener.Accept()
		if err != nil {
			continue
		}

		s.wg.Add(1)
		go s.handleClientConnection(socket)
	}
}

func (s *Server) handleClientConnection(socket net.Conn) {
	defer s.wg.Done()
	defer socket.Close()

	log := s.logger.With().Str("peer", socket.RemoteAddr().String()).Logger()
	log.Info().Msg("Accepted incoming connection")

	for {
		data := make([]byte, bufSize)
		bytesRead, err := socket.Read(data)
		switch err {
		case nil:
		case io.EOF:
			s.logger.Info().Msg("Connection closed")
			return
		default:
			s.logger.Err(err).Msg("Failed to read from socket")
			return
		}

		s.logger.Trace().
			Int("len", bytesRead).
			Hex("data", data[0:bytesRead]).
			Msg("Received data from socket")

		envelope, err := NewEnvelope(data[0:bytesRead])
		if err != nil {
			s.logger.Err(err).Msg("Invalid message")
			continue
		}

		s.logger.Info().Stringer("type", envelope.MessageType()).Msg("Received message")
		response, err := s.handler.handleMessage(context.TODO(), envelope)
		if err != nil {
			s.logger.Err(err).Msg("Failed to parse message")
		}

		if err := s.sendResponse(socket, response); err != nil {
			s.logger.Err(err).Msg("Failed to respond to message")
		}
	}
}

func (s *Server) sendResponse(client net.Conn, envelope Envelope) error {
	s.logger.Trace().
		Int("len", len(envelope)).
		Hex("data", envelope).
		Msg("Writing response to socket")

	if _, err := client.Write(envelope); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}
