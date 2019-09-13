package lg_webos

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/snabb/webostv"
)

func (b *Bind) monitorForegroundAppInfo(s webostv.ForegroundAppInfo) error {
	ctx := context.Background()
	sn := b.SerialNumber()

	if err := b.MQTTPublishAsync(ctx, b.config.TopicStateApplication.Format(sn), s.AppId); err != nil {
		return err
	}

	// TODO: cache
	if s.AppId == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
	}

	return b.MQTTPublishAsync(ctx, b.config.TopicStatePower.Format(sn), s.AppId != "")
}

func (b *Bind) monitorAudio(s webostv.AudioStatus) error {
	ctx := context.Background()
	sn := b.SerialNumber()

	if err := b.MQTTPublishAsync(ctx, b.config.TopicStateMute.Format(sn), s.Mute); err != nil {
		return err
	}

	return b.MQTTPublishAsync(ctx, b.config.TopicStateVolume.Format(sn), s.Volume)
}

func (b *Bind) monitorTvCurrentChannel(s webostv.TvCurrentChannel) error {
	return b.MQTTPublishAsync(context.Background(), b.config.TopicStateChannelNumber.Format(b.SerialNumber()), s.ChannelNumber)
}
