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

	switch v := err.Error(); {
	case
		strings.Contains(v, "use of closed network connection"),
		strings.Contains(v, "connection reset by peer"),
		strings.Contains(v, "bad file descriptor"):
		return err
	}

	return nil
}
