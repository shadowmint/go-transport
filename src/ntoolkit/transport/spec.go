package transport

// Spec is the configuration object to create a new transport object.
type Spec struct {

	// Port is the port to bind to.
	Port int

	// MaxThreads is the maxmimum number of concurrent connections to handle.
	MaxThreads int

	// Format defines the format to use for this connection.
	Format interface{}

	// Handler the the function to invoke on valid incoming connections.
	Handler func(*Api)
}

// FormatJSON is the format key to use for simple json TCP connections.
type FormatJSON struct{}
