// +build !linux !arm

package gpio

import (
	"fmt"
	"runtime"
)

func NewPin(_ int64, _ PinMode) (GPIOPin, error) {
	return nil, fmt.Errorf("Platform %s,%s isn't support", runtime.GOOS, runtime.GOARCH)
}
