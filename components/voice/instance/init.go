package instance

import (
	"github.com/kihamo/boggart/components/voice/internal"
	"github.com/kihamo/shadow"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}
