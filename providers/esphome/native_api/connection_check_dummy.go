// +build windows appengine

package native_api

import (
	"net"
)

func ConnectionCheck(c net.Conn) error {
	return nil
}
