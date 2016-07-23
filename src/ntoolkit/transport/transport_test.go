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
	trans := transport.New(func(api *transport.Api) {
		fmt.Printf("Actually invoked the handler!!!!! %v\n", api)
	}, nil)
	go func() {
		trans.Listen("127.0.0.1:0")
	}()

	time.Sleep(time.Second / 10)
	conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", trans.Port()))
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
