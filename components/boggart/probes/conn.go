package probes

import (
	"errors"
	"strings"
	"syscall"
)

func ConnErrorProbe(err error) error {
	if err == nil {
		return err
	}

	if errors.Is(err, syscall.EPIPE) {
		return err
	}

	if strings.Contains(err.Error(), "use of closed network connection") {
		return err
	}

	return nil
}
