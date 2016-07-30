package transport_test

import (
	"fmt"
	"net"
	"ntoolkit/assert"
	"ntoolkit/transport"
	"testing"
)

func TestNetworks(T *testing.T) {
	assert.Test(T, func(T *assert.T) {
		validate := func(ips []net.IP, err error) bool {
			if err != nil {
				return false
			}
			fmt.Printf("%v\n", ips)
			return true
		}
		T.Assert(validate(transport.Networks(true, true, true)))
		T.Assert(validate(transport.Networks(true, false, false)))
		T.Assert(validate(transport.Networks(false, true, false)))
		T.Assert(validate(transport.Networks(true, true, false)))
		T.Assert(validate(transport.Networks(false, false, true)))
	})
}
