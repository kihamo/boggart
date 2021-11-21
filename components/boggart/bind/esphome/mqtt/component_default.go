package mqtt

import (
	"fmt"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentDefault struct {
	*componentBase

	state atomic.Value
}

func NewComponentDefault(id string, t ComponentType, message mqtt.Message) *ComponentDefault {
	return &ComponentDefault{
		componentBase: newComponentBase(id, t, message),
	}
}

func (c *ComponentDefault) State() interface{} {
	return c.state.Load()
}

func (c *ComponentDefault) StateFormat() string {
	state := c.State()

	if state == nil {
		return ""
	}

	return fmt.Sprintf("%v", c.State())
}

func (c *ComponentDefault) SetState(message mqtt.Message) error {
	c.state.Store(message.String())

	return nil
}
