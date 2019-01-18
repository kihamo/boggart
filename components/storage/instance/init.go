package instance

import (
	"github.com/kihamo/boggart/components/storage/internal"
	"github.com/kihamo/shadow"
)

func NewComponent() shadow.Component {
	return &internal.Component{}
}

func init() {
	shadow.MustRegisterComponent(NewComponent())
}
