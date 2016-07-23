package transport

import (
	"fmt"
	"net"
)

type Transport struct {
	handler func(*Api)
	active  bool
}

// New creates a new transport instance with a handler.
func New(handler func(*Api)) *Transport {
	return &Transport{handler, false}
}

// Listen binds to a local port and begin listening for connections.
// The handler will be called in a goroutine for connections.
func (transport *Transport) Listen(port int, maxThreads int) error {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		fmt.Printf("%v\n", err)
		return err
	}

	go func() {
		fmt.Printf("Waiting for connections...\n")
		l.SetDeadline()

		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		fmt.Printf("Got a connection\n")
		fmt.Printf("%v\n", conn.RemoteAddr().Network())
		fmt.Printf("%v\n", conn.RemoteAddr().String())
	}()
	return nil
}

// Halt stops listening on the socket and halts the worker threads
func (transport *Transport) Halt() {
}
