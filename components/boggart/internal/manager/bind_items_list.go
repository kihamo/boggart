package manager

import (
	"github.com/kihamo/boggart/components/boggart"
)

type BindItemsList []boggart.BindItem

func (l BindItemsList) MarshalYAML() (interface{}, error) {
	return struct {
		Devices []boggart.BindItem
	}{
		Devices: l,
	}, nil
}
