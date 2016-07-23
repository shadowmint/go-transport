package transport

import "net"
import "ntoolkit/jsonbridge"

// API is the api available to transport event handlers
type API struct {
	Connection *net.Conn
	bridge     *jsonbridge.Bridge
}

// Read attempts to read the data segment into the given data object.
func (api *API) Read(data interface{}) error {
	return bridge.As(data)
}

// Write attempts to write the data given to the socket connection.
func (api *API) Write(data interface{}) error {
	return bridge.Write(data)
}
