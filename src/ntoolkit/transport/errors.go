package transport

// ErrBind is raised if binding the network socket failed.
type ErrBind struct{}

// ErrBadAddress is raised if a malformed TCP address is used.
type ErrBadAddress struct{}

// ErrNetworks is raised if an error occurred enumerating local network interfaces
type ErrNetworks struct{}
