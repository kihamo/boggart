package internal

import (
	"github.com/kihamo/boggart/components/boggart/internal/listeners"
)

func (c *Component) initListeners() {
	c.listenersManager.AddListener(listeners.NewLoggingListener(c.logger))
}
