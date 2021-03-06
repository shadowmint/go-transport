package transport_test

import (
	"fmt"
	"net"
	"ntoolkit/assert"
	"ntoolkit/transport"
	"testing"
	"time"
)

func TestRun(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		resolved := 0

		trans := transport.New(func(api *transport.API) {
			resolved += 1
			if resolved == 2 {
				api.Shutdown()
			}
		}, nil)

		go func() { trans.Listen("127.0.0.1:0") }()
		time.Sleep(time.Second / 10)

		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", trans.Port()))
		T.Assert(err == nil)

		if conn != nil {
			conn.Write([]byte("{\"hello\": \"world\"}\n"))
			conn.Write([]byte("{\"hello\": \"world\"}\n"))
		}

		trans.Wait()
		T.Assert(resolved == 2)
	})
}

func TestRemoteDisconnect(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		resolved := 0

		trans := transport.New(func(api *transport.API) {
			resolved += 1
			if resolved == 2 {
				api.Shutdown()
			}
		}, nil)
		go func() { trans.Listen("127.0.0.1:0") }()
		time.Sleep(time.Second / 10)

		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", trans.Port()))
		T.Assert(err == nil)

		if conn != nil {
			conn.Write([]byte("{\"hello\": \"world\"}\n"))
			conn.Write([]byte("{\"hello\": \"world\"}"))
			conn.Close()
		}

		trans.Wait()
		T.Assert(resolved == 2)
	})
}

func TestMultiConnection(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		resolved := 0

		trans := transport.New(func(api *transport.API) {
			resolved += 1
			if resolved == 2 {
				api.Shutdown()
			}
		}, nil)
		trans.Config.MaxThreads = 2

		go func() { trans.Listen("127.0.0.1:0") }()
		time.Sleep(time.Second / 10)

		conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", trans.Port()))
		T.Assert(err == nil)

		if conn != nil {
			conn.Write([]byte("{\"hello\": \"world\"}\n"))
		}

		conn, err = net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", trans.Port()))
		T.Assert(err == nil)

		if conn != nil {
			conn.Write([]byte("{\"hello\": \"world\"}\n"))
		}

		trans.Wait()
		T.Assert(resolved == 2)
	})
}
