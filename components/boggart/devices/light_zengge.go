package devices

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"github.com/vikstrous/zengge-lightcontrol/control"
	"github.com/vikstrous/zengge-lightcontrol/manage"
)

const (
	ZenggeLightLocalPort  = 5577
	ZenggeLightManagePort = 48899
)

type ZenggeLight struct {
	boggart.DeviceWithSerialNumber

	controller *control.Controller
	manager    *manage.Manager
}

func NewZenggeLight(transport control.Transport, manager *manage.Manager) *ZenggeLight {
	device := &ZenggeLight{
		controller: &control.Controller{
			Transport: transport,
		},
		manager: manager,
	}
	device.Init()
	device.SetDescription("Zengge light")

	return device
}

func (d *ZenggeLight) Types() []boggart.DeviceType {
	return []boggart.DeviceType{
		boggart.DeviceTypeLight,
	}
}

func (d *ZenggeLight) Ping(ctx context.Context) bool {
	state, err := d.controller.GetState()

	if err != nil {
		return false
	}

	return state.IsOn
}

func (d *ZenggeLight) Enable() error {
	err := d.controller.SetPower(true)
	if err != nil {
		return err
	}

	return d.DeviceBase.Enable()
}

func (d *ZenggeLight) Disable() error {
	err := d.controller.SetPower(false)
	if err != nil {
		return err
	}

	return d.DeviceBase.Disable()
}

func (d *ZenggeLight) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillStopTask(d.taskSerialNumber)
	taskSerialNumber.SetTimeout(time.Second * 5)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Minute)
	taskSerialNumber.SetName("device-light-zengge-serial-number")

	return []workers.Task{
		taskSerialNumber,
	}
}

func (d *ZenggeLight) taskSerialNumber(ctx context.Context) (interface{}, error, bool) {
	if !d.IsEnabled() {
		return nil, nil, false
	}

	response, err := d.manager.ReliableRequestReceive("HF-A11ASSISTHREAD")
	if err != nil {
		return nil, err, false
	}

	parts := strings.Split(response, ",")
	if len(parts) < 3 {
		return nil, fmt.Errorf("Unparsable response from lightbulb: %s", response), false
	}

	d.SetSerialNumber(parts[1])
	d.SetDescription("Zengge light type " + parts[2] + " with serial number " + parts[1])

	return nil, nil, true
}
