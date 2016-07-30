# Transport

TCP based data transport layer.

# Usage

    package main

    import (
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

    	trans.Listen("127.0.0.1:0")
    	logger.Printf("Listening on %d...\n", trans.Port())

    	trans.Wait()
    }
