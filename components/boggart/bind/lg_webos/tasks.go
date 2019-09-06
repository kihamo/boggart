package lg_webos

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.livenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.livenessInterval)
	taskLiveness.SetName("liveness-" + b.host)

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	client, err := b.Client()
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	_, err = client.Register(b.key)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, nil
	}

	if b.IsStatusOnline() {
		return nil, nil
	}

	if b.SerialNumber() == "" {
		deviceInfo, err := client.GetCurrentSWInformation()
		if err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
			return nil, err
		}

		b.SetSerialNumber(deviceInfo.DeviceId)
	}

	b.UpdateStatus(boggart.BindStatusOnline)

	// set tv subscribers
	// TODO: close if OFFLINE
	quit := make(chan struct{})

	go func() {
		state, err := client.ApplicationManagerGetForegroundAppInfo()
		if err == nil {
			err = b.monitorForegroundAppInfo(state)
		}

		if err != nil {
			b.Logger().Warn("Failed get current app info", "error", err.Error())
		}

		err = client.ApplicationManagerMonitorForegroundAppInfo(b.monitorForegroundAppInfo, quit)
		if err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
		}

		// TODO: send to quit ???
	}()

	go func() {
		state, err := client.AudioGetStatus()
		if err == nil {
			err = b.monitorAudio(state)
		}

		if err != nil {
			b.Logger().Warn("Failed get current audio status", "error", err.Error())
		}

		err = client.AudioMonitorStatus(b.monitorAudio, quit)
		if err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
		}
	}()

	go func() {
		state, err := client.TvGetCurrentChannel()
		if err == nil {
			err = b.monitorTvCurrentChannel(state)
		}

		if err != nil {
			b.Logger().Warn("Failed get current tv channel", "error", err.Error())
		}

		err = client.TvMonitorCurrentChannel(b.monitorTvCurrentChannel, quit)
		if err != nil {
			b.UpdateStatus(boggart.BindStatusOffline)
		}
	}()

	return nil, nil
}
