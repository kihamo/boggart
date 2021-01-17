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

func NewComponentDefault(id string, t ComponentType) *ComponentDefault {
	return &ComponentDefault{
		componentBase: newComponentBase(id, t),
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
