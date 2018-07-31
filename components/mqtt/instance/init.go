package instance

import (
	"github.com/kihamo/boggart/components/mqtt/internal"
	"github.com/kihamo/shadow"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}
