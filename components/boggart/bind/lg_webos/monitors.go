package lg_webos

import (
	"context"

	"github.com/snabb/webostv"
	"go.uber.org/multierr"
)

func (b *Bind) monitorForegroundAppInfo(s webostv.ForegroundAppInfo) (err error) {
	ctx := context.Background()
	sn := b.Meta().SerialNumber()

	b.power.Set(s.AppId != "")

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateApplication.Format(sn), s.AppId); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStatePower.Format(sn), s.AppId != ""); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) monitorAudio(s webostv.AudioStatus) error {
	ctx := context.Background()
	sn := b.Meta().SerialNumber()

	if err := b.MQTT().PublishAsync(ctx, b.config.TopicStateMute.Format(sn), s.Mute); err != nil {
		return err
	}

	return b.MQTT().PublishAsync(ctx, b.config.TopicStateVolume.Format(sn), s.Volume)
}

func (b *Bind) monitorTvCurrentChannel(s webostv.TvCurrentChannel) error {
	return b.MQTT().PublishAsync(context.Background(), b.config.TopicStateChannelNumber.Format(b.Meta().SerialNumber()), s.ChannelNumber)
}
