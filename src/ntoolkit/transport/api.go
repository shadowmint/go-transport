package transport

import "net"

type Api struct {
	Connection *net.Conn
}

// Read attempts to read the data segment into the given data object.
func (api *Api) Read(data interface{}) error {
	return nil
}

// Write attempts to write the data given to the socket connection.
func (api *Api) Write(data interface{}) error {
	return nil
}
