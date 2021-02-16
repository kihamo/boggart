package fcm

import (
	"github.com/kihamo/boggart/components/boggart"
)

/*
~/openhab-cloud$ mongo 127.0.0.1:27017/openhab --eval "db.userdevices.find()"
*/

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{}
}
