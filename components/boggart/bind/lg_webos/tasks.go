package webos

import (
	"context"
	"errors"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("serial-number").
			WithHandler(
				b.Workers().WrapTaskIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSerialNumber),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
					1,
				),
			),
	}
}

func (b *Bind) taskSerialNumber(ctx context.Context) error {
	client := b.Client()
	if client == nil {
		return errors.New("client isn't init")
	}

	deviceInfo, err := client.GetCurrentSWInformation()
	if err != nil {
		return err
	}

	b.Meta().SetSerialNumber(deviceInfo.DeviceId)
	b.Meta().SetMACAsString(deviceInfo.DeviceId)

	// set tv subscribers
	go func() {
		var err error

		// current state
		state, err := client.ApplicationManagerGetForegroundAppInfo()
		if err == nil {
			err = b.monitorForegroundAppInfo(state)
		}

		if err != nil {
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
		state, err := client.AudioGetStatus()
		if err == nil {
			err = b.monitorAudio(state)
		}

		if err != nil {
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
		state, err := client.TvGetCurrentChannel()
		if err == nil {
			err = b.monitorTvCurrentChannel(state)
		}

		if err != nil {
			b.Logger().Warn("Failed get current tv channel", "error", err.Error())
		}

		// subscriber
		if err = client.TvMonitorCurrentChannel(b.monitorTvCurrentChannel, b.quitMonitors); err != nil {
			b.Logger().Errorf("Init channel monitor failed %v", err)
		}
	}()

	return nil
}
