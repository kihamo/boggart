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
	taskLiveness.SetName("bind-lg-webos-liveness-" + b.host)

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	client, err := b.Client()
	if err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	_, err = client.Register(b.key)
	if err != nil {
		b.UpdateStatus(boggart.DeviceStatusOffline)
		return nil, err
	}

	if b.Status() == boggart.DeviceStatusOnline {
		return nil, nil
	}

	if b.SerialNumber() == "" {
		deviceInfo, err := client.GetCurrentSWInformation()
		if err != nil {
			b.UpdateStatus(boggart.DeviceStatusOffline)
			return nil, err
		}

		b.SetSerialNumber(deviceInfo.DeviceId)
	}

	b.UpdateStatus(boggart.DeviceStatusOnline)

	// set tv subscribers
	// TODO: close if OFFLINE
	quit := make(chan struct{})

	go func() {
		state, err := client.ApplicationManagerGetForegroundAppInfo()
		if err == nil {
			b.monitorForegroundAppInfo(state)
		}

		client.ApplicationManagerMonitorForegroundAppInfo(b.monitorForegroundAppInfo, quit)
	}()

	go func() {
		state, err := client.AudioGetStatus()
		if err == nil {
			b.monitorAudio(state)
		}

		client.AudioMonitorStatus(b.monitorAudio, quit)
	}()

	go func() {
		state, err := client.TvGetCurrentChannel()
		if err == nil {
			b.monitorTvCurrentChannel(state)
		}

		client.TvMonitorCurrentChannel(b.monitorTvCurrentChannel, quit)
	}()

	return nil, nil
}
