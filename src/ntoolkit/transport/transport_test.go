package transport_test

import (
	"fmt"
	"net"
	"ntoolkit/transport"
	"testing"
	"time"
)

func TestRun(T *testing.T) {

	// Open outgoing connection
	go func() {
		trans := transport.New(func(api *transport.Api) {

		})
		trans.Listen(5000, 1)
	}()

	time.Sleep(time.Second / 10)
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	if conn != nil {
		fmt.Printf("connected from test")
		conn.Write([]byte("{\"hello\": \"world\"}"))

		time.Sleep(time.Second / 10)
		conn.Close()
	}
}
