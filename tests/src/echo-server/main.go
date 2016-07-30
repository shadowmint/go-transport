package main

import (
	"fmt"
	"log"
	"ntoolkit/transport"
	"os"
)

// Basic message type
type message struct {
	Message string
}

func main() {

	logger := log.New(os.Stdout, "EchoServer: ", log.Ldate|log.Ltime|log.Lshortfile)

	config := transport.Config{
		MaxThreads:    2,
		AcceptTimeout: 100,
		ReadTimeout:   1000,
		Logger:        logger,
	}

	// Handle messages
	trans := transport.New(func(api *transport.API) {
		var msg message
		if err := api.Read(&msg); err == nil {
			api.Logger.Printf("Got message: %s\n", msg.Message)
			api.Write(msg)
			if msg.Message == "END" {
				api.Logger.Printf("Got END request. Closing connection...")
				api.Close()
			} else if msg.Message == "EXIT" {
				api.Logger.Printf("Got EXIT request. Closing server...")
				api.Shutdown()
			}
		} else {
			api.Logger.Printf("Unknown message format: %s\n", api.Raw())
		}
	}, &config)

	// Find a local loopback to bind on tcp4
	networks, err := transport.Networks(true, true, false)
	if err != nil {
		logger.Printf("Failed to start: %v\n", err)
	} else if len(networks) == 0 {
		logger.Printf("Failed to start: No local network interfaces found.\n")
	}

	// Start listening
	host := fmt.Sprintf("%s:0", networks[0])
	trans.Listen(host)
	logger.Printf("Listening on %s:%d...\n", networks[0], trans.Port())

	// Wait for end~
	trans.Wait()
}
