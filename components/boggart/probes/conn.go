package probes

import (
	"errors"
	"syscall"
)

func ConnErrorProbe(err error) error {
	if errors.Is(err, syscall.EPIPE) {
		return err
	}

	return nil
}
