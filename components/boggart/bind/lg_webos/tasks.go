package lg_webos

import (
	"context"
	"errors"
	"time"

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

	client := b.Client()
	if client == nil {
		return nil, errors.New("client isn't init")
	}

	deviceInfo, err := client.GetCurrentSWInformation()
	if err != nil {
		return nil, err
	}

	b.Meta().SetSerialNumber(deviceInfo.DeviceId)

	// set tv subscribers
	go func() {
		var err error

		// current state
		if state, err := client.ApplicationManagerGetForegroundAppInfo(); err == nil {
			err = b.monitorForegroundAppInfo(state)
		} else {
			b.Logger().Warn("Failed get current app info", "error", err.Error())
		}

		// subscriber
		if err = client.ApplicationManagerMonitorForegroundAppInfo(b.monitorForegroundAppInfo, b.quitMonitors); err != nil {
			b.Logger().Errorf("Init application manager monitor failed %v", err)
		}
	}()

	go func() {
		var err error

		// current state
		if state, err := client.AudioGetStatus(); err == nil {
			err = b.monitorAudio(state)
		} else {
			b.Logger().Warn("Failed get current audio status", "error", err.Error())
		}

		// subscriber
		if err = client.AudioMonitorStatus(b.monitorAudio, b.quitMonitors); err != nil {
			b.Logger().Errorf("Init audio monitor failed %v", err)
		}
	}()

	go func() {
		var err error

		// current state
		if state, err := client.TvGetCurrentChannel(); err == nil {
			err = b.monitorTvCurrentChannel(state)
		} else {
			b.Logger().Warn("Failed get current tv channel", "error", err.Error())
		}

		// subscriber
		if err = client.TvMonitorCurrentChannel(b.monitorTvCurrentChannel, b.quitMonitors); err != nil {
			b.Logger().Errorf("Init channel monitor failed %v", err)
		}
	}()

	return nil, nil
}
