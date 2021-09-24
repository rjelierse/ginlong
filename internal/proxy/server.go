package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/rjelierse/ginlong/internal/messages"
)

// Proxy sits in between a local data logging stick and the upstream monitoring system.
type Proxy struct {
	localAddress    string
	upstreamAddress string

	messages     chan messages.Envelope
	measurements chan messages.Measurement
	shutdown     chan struct{}
}

// New constructs a new Proxy server.
func New(localAddress, upstreamAddress string) *Proxy {
	return &Proxy{
		localAddress:    localAddress,
		upstreamAddress: upstreamAddress,
		messages:        make(chan messages.Envelope, 100),
		measurements:    make(chan messages.Measurement, 100),
		shutdown:        make(chan struct{}),
	}
}

// Listen starts the server in listening mode.
//
// This function blocks until Shutdown is called.
func (p *Proxy) Listen(ctx context.Context) error {
	lc := net.ListenConfig{}
	l, err := lc.Listen(ctx, "tcp", p.localAddress)
	if err != nil {
		return err
	}
	log.Println("Listening for connections:", l.Addr())

	go p.handleConnections(l)

	p.processMessages()

	return nil
}

// Shutdown signals the Proxy should terminate.
func (p *Proxy) Shutdown() {
	close(p.shutdown)
}

// Measurements returns a channel for receiving the measurement data.
func (p *Proxy) Measurements() <-chan messages.Measurement {
	return p.measurements
}

func (p *Proxy) handleConnections(l net.Listener) {
	wg := sync.WaitGroup{}

	conns := make(chan *Connection)
	go func() {
		for {
			c, err := l.Accept()
			if err != nil {
				log.Println("Failed to accept connection:", err)
				return
			}
			upstream, err := net.Dial("tcp", p.upstreamAddress)
			if err != nil {
				log.Println("Failed to connect to upstream:", err)
				return
			}
			conns <- &Connection{
				client:   c,
				upstream: upstream,
				shutdown: p.shutdown,
				messages: p.messages,
				errors:   make(chan error),
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-p.shutdown:
				log.Println("Shutting down listener:", l.Close())
				return
			case c := <-conns:
				log.Printf("[%s] Connected.\n", c.client.RemoteAddr())

				wg.Add(1)
				go func() {
					defer wg.Done()
					c.run()
				}()
			}
		}
	}()

	// Drain connections
	wg.Wait()

	// Close messages channel
	close(p.messages)
}

func (p *Proxy) processMessages() {
	for msg := range p.messages {
		p.handleMessage(msg)
	}
}

func (p *Proxy) handleMessage(msg messages.Envelope) {
	fmt.Printf("Envelope%+v\n", msg)

	if msg.MessageType != messages.MessageTypeMeasurement {
		return
	}

	m, err := messages.NewMeasurement(msg.Message)
	if err != nil {
		log.Println("Invalid measurement:", err)
		return
	}

	p.measurements <- m
}
