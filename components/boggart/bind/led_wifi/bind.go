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
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind

	bulb *wifiled.Bulb
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.bulb = wifiled.NewBulb(b.config().Address)

	return nil
}

func (b *Bind) State(ctx context.Context) error {
	state, err := b.bulb.State(ctx)
	if err != nil {
		return err
	}

	b.Meta().SetSerialNumber(strconv.FormatUint(uint64(state.DeviceName), 10))

	var mqttErr error
	cfg := b.config()

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStatePower.Format(cfg.Address), state.Power); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateMode.Format(cfg.Address), state.Mode); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateSpeed.Format(cfg.Address), state.Speed); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	// in HEX
	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateColor.Format(cfg.Address), state.Color.String()); e != nil {
		mqttErr = multierror.Append(mqttErr, e)
	}

	// in HSV
	h, s, v := state.Color.HSV()
	if e := b.MQTT().PublishAsync(ctx, cfg.TopicStateColorHSV.Format(cfg.Address), fmt.Sprintf("%d,%.2f,%.2f", h, s, v)); e != nil {
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
