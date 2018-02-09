package internal

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/devices"
	"github.com/kihamo/boggart/components/boggart/providers/hikvision"
)

func (c *Component) initCameras() {
	isapi := hikvision.NewISAPI(
		c.config.String(boggart.ConfigHikvisionHallHost),
		c.config.Int64(boggart.ConfigHikvisionHallPort),
		c.config.String(boggart.ConfigHikvisionHallUsername),
		c.config.String(boggart.ConfigHikvisionHallPassword))

	device := devices.NewHikVisionCamera(isapi, c.config.Uint64(boggart.ConfigHikvisionHallStreamingChannel))
	c.devices.Register(boggart.DeviceCameraHallID, device)

	isapi = hikvision.NewISAPI(
		c.config.String(boggart.ConfigHikvisionStreetHost),
		c.config.Int64(boggart.ConfigHikvisionStreetPort),
		c.config.String(boggart.ConfigHikvisionStreetUsername),
		c.config.String(boggart.ConfigHikvisionStreetPassword))

	device = devices.NewHikVisionCamera(isapi, c.config.Uint64(boggart.ConfigHikvisionStreetStreamingChannel))
	c.devices.Register(boggart.DeviceCameraStreetID, device)
}
