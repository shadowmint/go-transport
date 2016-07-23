package transport

import "net"
import "ntoolkit/jsonbridge"

type Api struct {
	bridge     *jsonbridge.Bridge
	Connection *net.Conn
}

// Read attempts to read the data segment into the given data object.
func (api *Api) Read(data interface{}) error {
	return api.bridge.As(data)
}

// Write attempts to write the data given to the socket connection.
func (api *Api) Write(data interface{}) error {
	return api.bridge.Write(data)
}
