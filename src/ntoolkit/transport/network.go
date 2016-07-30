package transport

import (
	"net"
	"ntoolkit/errors"
	"strings"
)

// Networks returns a list of local network interfaces that could potentially
// have a service bound to them.
func Networks(loopback bool, tcp4 bool, tcp6 bool) ([]net.IP, error) {
	rtn := make([]net.IP, 0, 5)
	ifaces, err := net.Interfaces()
	if err != nil {
		return rtn, errors.Fail(ErrNetworks{}, err, "Unable to get network interface list")
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err == nil {
			for _, addr := range addrs {
				switch v := addr.(type) {
				case *net.IPNet:
					if CheckIPValid(v, loopback, tcp4, tcp6) {
						rtn = append(rtn, v.IP)
					}
				}
			}
		}
	}
	return rtn, nil
}

// CheckIPValid returns true or false based on if a net.IPNet matches the
// given set of requirements from isLocal, isTCP4 and isTCP6.
func CheckIPValid(net *net.IPNet, isLocal bool, isTCP4 bool, isTCP6 bool) bool {
	aLocal, aTCP4, aTCP6 := CheckIPType(net)
	return ((isLocal == aLocal) && (isTCP4 == aTCP4)) || ((isLocal == aLocal) && (isTCP6 == aTCP6))
}

// CheckIPType returns a tuple of (loopback, tcp4, tcp6) for an net.IPNet
func CheckIPType(net *net.IPNet) (bool, bool, bool) {
	var isLocal = net.IP.IsLoopback() || strings.HasSuffix(net.String(), "::1/64")
	var isTCP4 = net.IP.To4() != nil
	var isTCP6 = !isTCP4 && net.IP.To16() != nil
	return isLocal, isTCP4, isTCP6
}
