package transport

import (
	"fmt"
	"log"
	"net"
	"ntoolkit/errors"
	"ntoolkit/jsonbridge"
	"ntoolkit/threadpool"
	"time"
)

// Config is the set of named config values for a Transport instance.
type Config struct {

	// MaxThreads is the maximum number of concurrent active handlers.
	MaxThreads int

	// AcceptTimeout is the maximum blocking length of Accept() calls.
	AcceptTimeout int

	// ReadTimeout is the maximum blocking length of read requests.
	ReadTimeout int

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
			ReadTimeout:   100,
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
	transport.pool = threadpool.New()
	transport.pool.MaxThreads = transport.Config.MaxThreads

	// Handle requests
	go func() {
		for {
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
				api := &Api{Connection: &conn}
				err := transport.pool.Run(func() {

					// Setup
					bridge := jsonbridge.New(conn, conn)
					bridge.Timeout = transport.Config.ReadTimeout
					api.bridge = bridge

					// Read content on the connection until it closes and push to the handler
					// when a completed token is ready.
					for transport.active {
						bridge.Read()
						for bridge.Len() > 0 {
							fmt.Printf("Got a token to handle!\n")
							bridge.Next()
							transport.handler(api)
						}
					}
				})
				if errors.Is(err, threadpool.ErrBusy{}) {
					transport.logWarning("Failed to handle incoming connect; no available handlers")
				} else {
					transport.logError("Error handling connection", err)
				}
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

// Log some error message
func (transport *Transport) logError(message string, err error) {
	if transport.Config.Logger != nil {
		fmt.Printf("%v %v", message, err)
		transport.Config.Logger.Print(err)
		transport.Config.Logger.Print(message)
	}
}

// Log some warning message
func (transport *Transport) logWarning(message string) {
	if transport.Config.Logger != nil {
		fmt.Printf("%v", message)
		transport.Config.Logger.Print(message)
	}
}
