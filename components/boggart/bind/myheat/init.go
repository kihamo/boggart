package myheat

import (
	"net/url"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/bind/myheat/cloud"
	"github.com/kihamo/boggart/components/boggart/bind/myheat/device"
	"github.com/kihamo/boggart/providers/myheat/cloud/client"
)

func init() {
	link, err := url.Parse(client.DefaultHost)
	if err == nil {
		link.Scheme = client.DefaultSchemes[0]
	}

	boggart.RegisterBindType("myheat", device.Type{}, "myheat:device", "myheat:smart2", "myheat:device:smart2")
	boggart.RegisterBindType("myheat:cloud", cloud.Type{
		Link: link,
	})
}
