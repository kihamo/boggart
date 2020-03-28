package ledwifi

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/providers/wifiled"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind

	config *Config
	bulb   *wifiled.Bulb
}

func (b *Bind) State(ctx context.Context) error {
	state, err := b.bulb.State(ctx)
	if err != nil {
		return err
	}

	b.Meta().SetSerialNumber(strconv.FormatUint(uint64(state.DeviceName), 10))
	var mqttErr error

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStatePower, state.Power); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateMode, state.Mode); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateSpeed, state.Speed); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	// in HEX
	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateColor, state.Color.String()); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	// in HSV
	h, s, v := state.Color.HSV()
	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateColorHSV, fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	if mqttErr != nil {
		b.Logger().Error(mqttErr.Error())
	}

	return nil
}

func (b *Bind) On(ctx context.Context) error {
	err := b.bulb.PowerOn(ctx)
	if err == nil {
		err = b.State(ctx)
	}

	return err
}

func (b *Bind) Off(ctx context.Context) error {
	err := b.bulb.PowerOff(ctx)
	if err == nil {
		err = b.State(ctx)
	}

	return err
}
