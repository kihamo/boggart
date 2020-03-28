// +build windows appengine

package nativeapi

import (
	"net"
)

func ConnectionCheck(c net.Conn) error {
	return nil
}
