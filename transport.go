package gosc

import (
	"net"
)

// Transport interface describes the transportation used by the Client and
// Server.
type Transport interface {
	Send(pack Package, addr net.Addr) error
	Receive() (pack Package, from net.Addr, err error)
}
