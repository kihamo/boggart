//go:build windows || appengine
// +build windows appengine

package connection

import (
	"net"
)

func Check(c net.Conn) error {
	return nil
}
