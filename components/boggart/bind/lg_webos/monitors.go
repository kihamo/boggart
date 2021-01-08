package webos

import (
	"context"

	"github.com/snabb/webostv"
	"go.uber.org/multierr"
)

func (b *Bind) monitorForegroundAppInfo(s webostv.ForegroundAppInfo) (err error) {
	ctx := context.Background()
	mac := b.Meta().MACAsString()

	b.power.Set(s.AppId != "")

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStateApplication.Format(mac), s.AppId); e != nil {
		err = multierr.Append(err, e)
	}

	if e := b.MQTT().PublishAsync(ctx, b.config.TopicStatePower.Format(mac), s.AppId != ""); e != nil {
		err = multierr.Append(err, e)
	}

	return err
}

func (b *Bind) monitorAudio(s webostv.AudioStatus) error {
	ctx := context.Background()
	mac := b.Meta().MACAsString()

	if err := b.MQTT().PublishAsync(ctx, b.config.TopicStateMute.Format(mac), s.Mute); err != nil {
		return err
	}

	return b.MQTT().PublishAsync(ctx, b.config.TopicStateVolume.Format(mac), s.Volume)
}

func (b *Bind) monitorTvCurrentChannel(s webostv.TvCurrentChannel) error {
	return b.MQTT().PublishAsync(context.Background(), b.config.TopicStateChannelNumber.Format(b.Meta().MACAsString()), s.ChannelNumber)
}
