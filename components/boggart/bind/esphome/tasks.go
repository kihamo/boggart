package esphome

import (
	"context"

	"github.com/golang/protobuf/proto"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/providers/esphome/native_api"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("liveness-" + b.config.Address)

	taskUpdater := task.NewFunctionTask(b.taskUpdated)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater-" + b.config.Address)

	return []workers.Task{
		taskLiveness,
		taskUpdater,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	info, err := b.provider.DeviceInfo(ctx)
	if err != nil {
		b.UpdateStatus(boggart.BindStatusOffline)
		return nil, err
	}

	if sn := b.SerialNumber(); sn == "" {
		b.SetSerialNumber(info.MacAddress)
	}

	b.UpdateStatus(boggart.BindStatusOnline)
	return nil, nil
}

func (b *Bind) taskUpdated(ctx context.Context) (interface{}, error) {
	sn := b.SerialNumber()
	if sn == "" {
		return nil, nil
	}

	list, err := b.provider.ListEntities(ctx)
	if err != nil {
		return nil, err
	}

	states, err := b.States(ctx, list)
	if err != nil {
		return nil, err
	}

	entities := make(map[uint32]native_api.MessageEntity)
	for _, message := range list {
		if e, ok := message.(native_api.MessageEntity); ok {
			entities[e.GetKey()] = e
		}
	}

	for _, message := range states {
		s, ok := message.(native_api.MessageState)
		if !ok {
			continue
		}

		var entity native_api.MessageEntity
		entity, ok = entities[s.GetKey()]
		if !ok {
			continue
		}

		state, _, e := native_api.State(entity.(proto.Message), message)
		if e == nil {
			if e = b.MQTTPublishAsync(ctx, b.config.TopicState.Format(sn, entity.GetObjectId()), state); e != nil {
				err = multierr.Append(err, e)
			}
		} else {
			err = multierr.Append(err, e)
		}
	}

	return nil, nil
}
