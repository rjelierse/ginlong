package proxy

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/rjelierse/ginlong/internal/messages"
)

// Connection represents an incoming connection from a logging stick.
// Each connection has its own connection to the upstream monitoring system.
type Connection struct {
	client   net.Conn
	upstream net.Conn
	messages chan<- messages.Envelope
	shutdown <-chan struct{}
	errors   chan error
}

func (c *Connection) run() {
	defer c.client.Close()
	defer c.upstream.Close()

	go c.handleMessages()

	for {
		select {
		case <-c.shutdown:
			log.Printf("[%s] Closing connection.", c.client.RemoteAddr())
			return
		case err := <-c.errors:
			log.Printf("[%s] Error: %+v\n", c.client.RemoteAddr(), err)
			return
		}
	}
}

func (c *Connection) handleMessages() {
	for {
		msg, err := relayMessage(c.client, c.upstream)
		if err != nil {
			c.errors <- fmt.Errorf("failed to relay message to upstream: %w", err)
			break
		}

		c.messages <- messages.NewEnvelope(msg)

		res, err := relayMessage(c.upstream, c.client)
		if err != nil {
			c.errors <- fmt.Errorf("failed to relay response from upstream: %w", err)
			break
		}

		c.messages <- messages.NewEnvelope(res)
	}
}

func relayMessage(src io.Reader, dst io.Writer) ([]byte, error) {
	buf := make([]byte, 4096)
	bytesRecv, err := src.Read(buf)
	if err != nil {
		return nil, err
	}
	msg := buf[:bytesRecv]
	bytesSent, err := dst.Write(msg)
	if err != nil {
		return nil, err
	}
	if bytesRecv != bytesSent {
		return nil, fmt.Errorf("bytes received does not match bytes sent: %d v %d", bytesRecv, bytesSent)
	}
	return msg, nil
}
