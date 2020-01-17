package google_home

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/providers/google/home"
	"github.com/kihamo/boggart/providers/google/home/client/info"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskSerialNumber := task.NewFunctionTillSuccessTask(b.taskSerialNumber)
	taskSerialNumber.SetRepeats(-1)
	taskSerialNumber.SetRepeatInterval(time.Second * 30)
	taskSerialNumber.SetName("serial-number")

	return []workers.Task{
		taskSerialNumber,
	}
}

func (b *Bind) taskSerialNumber(ctx context.Context) (interface{}, error) {
	if !b.Meta().IsStatusOnline() {
		return nil, errors.New("bind isn't online")
	}

	response, err := b.provider.Info.GetEurekaInfo(info.NewGetEurekaInfoParams().
		WithOptions(home.EurekaInfoOptionDetail.Value()).
		WithParams(home.EurekaInfoParamDeviceInfo.Value()))

	if err != nil {
		return nil, err
	}

	if response.Payload == nil || response.Payload.MacAddress == "" {
		return nil, errors.New("MAC address not found")
	}

	b.Meta().SetSerialNumber(response.Payload.MacAddress)
	return nil, nil
}
