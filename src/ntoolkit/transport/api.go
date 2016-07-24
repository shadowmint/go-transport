package transport

import (
	"log"
	"net"
	"ntoolkit/jsonbridge"
)

// API is the api available to transport event handlers
type API struct {
	Connection net.Conn
	Logger     *log.Logger
	transport  *Transport
	bridge     *jsonbridge.Bridge
	active     bool
}

// Raw returns the current raw data block.
func (api *API) Raw() string {
	return api.bridge.Raw()
}

// Read attempts to read the data segment into the given data object.
func (api *API) Read(data interface{}) error {
	return api.bridge.As(data)
}

// Write attempts to write the data given to the socket connection.
func (api *API) Write(data interface{}) error {
	return api.bridge.Write(data)
}

// Context returns the context for this request, if any.
// Remember connections happen in their own threads, be sure to lock
// before modifying the context object if it is not thread safe.
func (api *API) Context() interface{} {
	return api.transport.Context
}

// Close the connection
func (api *API) Close() {
	api.active = false
	api.Connection.Close()
}

// Shutdown the entire server
func (api *API) Shutdown() {
	api.transport.Halt()
}
