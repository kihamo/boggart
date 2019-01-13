package broadlink

import (
	"time"
)

const (
	RMCaptureDuration = time.Second * 15
)

type ConfigRM struct {
	IP              string `valid:"ip,required"`
	MAC             string `valid:"mac,required"`
	Model           string `valid:"in(rm3mini|rm2proplus),required"`
	CaptureDuration time.Duration
}

func (t TypeRM) Config() interface{} {
	return &ConfigRM{
		CaptureDuration: RMCaptureDuration,
	}
}
