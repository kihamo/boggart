package devices

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/samsung/tv"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

// http://www.zoobab.com/samsung-tv-ue40d6200

/*
nmap -p 1-65535 192.168.88.169

Starting Nmap 7.40 ( https://nmap.org ) at 2018-02-26 23:47 MSK
Nmap scan report for 192.168.88.169
Host is up (0.0012s latency).
Not shown: 65517 closed ports
PORT      STATE SERVICE
7676/tcp  open  imqbrokerd
7678/tcp  open  unknown
8001/tcp  open  vcom-tunnel
8002/tcp  open  teradataordbms
8080/tcp  open  http-proxy
8187/tcp  open  unknown
9012/tcp  open  unknown
9119/tcp  open  mxit
9197/tcp  open  unknown
15500/tcp open  unknown
20490/tcp open  unknown
32768/tcp open  filenet-tms
32769/tcp open  filenet-rpc
32770/tcp open  sometimes-rpc3
32771/tcp open  sometimes-rpc5
41647/tcp open  unknown
45900/tcp open  unknown
48711/tcp open  unknown
*/

type SamsungTV struct {
	boggart.DeviceWithSerialNumber

	api *tv.ApiV2
}

func NewSamsungTV(api *tv.ApiV2) *SamsungTV {
	device := &SamsungTV{
		api: api,
	}
	device.Init()
	device.SetDescription("Samsung TV")

	return device
}

func (d *SamsungTV) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeTV,
	}
}

func (d *SamsungTV) Ping(ctx context.Context) bool {
	_, err := d.api.Device(ctx)
	return err == nil
}

func (d *SamsungTV) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-tv-samsung-serial-number")

	return []workers.Task{
		taskSerialNumber,
	}
}

func (d *SamsungTV) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	deviceInfo, err := d.api.Device(ctx)
	if err != nil {
		return nil, err, false
	}

	if deviceInfo.ID == "" {
		return nil, errors.New("Device returns empty serial number"), false
	}

	d.SetSerialNumber(deviceInfo.ID)
	d.SetDescription("Samsung TV with serial number " + deviceInfo.ID)

	return nil, nil, true
}
