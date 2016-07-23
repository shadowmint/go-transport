package transport

import (
	"fmt"
	"log"
	"net"
	"ntoolkit/errors"
	"ntoolkit/threadpool"
	"time"
)

// Config is the set of named config values for a Transport instance.
type Config struct {

	// MaxThreads is the maximum number of concurrent active handlers.
	MaxThreads int

	// AcceptTimeout is the maximum blocking length of Accept() calls.
	AcceptTimeout int

	// The logger to use with this transport
	Logger *log.Logger
}

// Transport is a local raw TCP listener for JSON objects.
type Transport struct {
	Config  *Config
	handler func(*Api)
	port    int
	active  bool
	pool    *threadpool.ThreadPool
}

// New creates a new transport instance with a handler.
// If no config object is passed, defaults are used.
func New(handler func(*Api), config *Config) *Transport {
	if config == nil {
		config = &Config{
			MaxThreads:    1,
			AcceptTimeout: 100,
			Logger:        nil}
	}
	return &Transport{
		Config:  config,
		handler: handler,
		active:  false}
}

// Listen resolves the addr string using net.ResolveTCPAddr
// and binds to it to listen for incoming connections.
// The handler will be called in a goroutine for connections.
func (transport *Transport) Listen(addr string) error {

	// Resolve address
	binding, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return errors.Fail(ErrBadAddress{}, err, "Unable to resolve TCP address to listen on")
	}

	// Create a listener
	l, err := net.ListenTCP("tcp", binding)
	if err != nil {
		return errors.Fail(ErrBind{}, err, "Unable to bind socket")
	}

	// Prepare to handle requests
	transport.port = l.Addr().(*net.TCPAddr).Port
	transport.active = true

	// Handle requests
	go func() {
		for true {
			// Try to handle the next connection
			l.SetDeadline(time.Now().Add(time.Millisecond * time.Duration(transport.Config.AcceptTimeout)))
			conn, err := l.Accept()
			if err != nil {
				if e, ok := err.(net.Error); ok && e.Timeout() {
					// On timeout, check if we need to stop
					if !transport.active {
						break
					}
				} else {
					// Some real error occurred
					transport.logError("Halting transport", err)
					break
				}
			} else {
				// Handle the connection?
				fmt.Printf("Got a connection\n")
				fmt.Printf("%v\n", conn)
				fmt.Printf("%v\n", err)

				pool := threadpool.New()
				pool.MaxThreads = 2

				value := 0

				T.Assert(pool.Run(func() { value += 1 }) == nil)
				T.Assert(pool.Run(func() { value += 1 }) == nil)
				err := pool.Run(func() { value += 1 })

				T.Assert(err != nil)
				T.Assert(errors.Is(err, threadpool.ErrBusy{}))

				pool.Wait()
				T.Assert(value == 2)
			}
		}

		// Restore state of the transport
		transport.active = false
	}()

	return nil
}

// Port returns the port that is currently bound for the listener, to support
// '127.0.0.1:0' style connection strings where the port is automatically assigned.
func (transport *Transport) Port() int {
	return transport.port
}

// Halt stops listening on the socket and halts the worker threads
func (transport *Transport) Halt() {
}

// Log some message
func (transport *Transport) logError(message string, err error) {
	if transport.Config.Logger != nil {
		fmt.Printf("%v %v", message, err)
		transport.Config.Logger.Print(err)
		transport.Config.Logger.Print(message)
	}
}
