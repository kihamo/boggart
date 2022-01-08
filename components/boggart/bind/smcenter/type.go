package smcenter

import (
	"net/url"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct {
	Link            *url.URL
	BaseURL         string
	BillContentType string
}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{}
}
