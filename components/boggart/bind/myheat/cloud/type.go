package cloud

import (
	"net/url"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	Link *url.URL
}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{}
}
