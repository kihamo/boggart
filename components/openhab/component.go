package openhab

import (
	"github.com/kihamo/boggart/components/openhab/client"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	Client() *client.OpenHABREST
}
