package lg_webos

import (
	"context"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/snabb/webostv"
)

func (b *Bind) monitorForegroundAppInfo(s webostv.ForegroundAppInfo) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(b.SerialNumber())

	if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateApplication.Format(sn), s.AppId); err != nil {
		return err
	}

	// TODO: cache
	if s.AppId == "" {
		b.UpdateStatus(boggart.BindStatusOffline)
	}

	return b.MQTTPublishAsync(ctx, MQTTPublishTopicStatePower.Format(sn), s.AppId != "")
}

func (b *Bind) monitorAudio(s webostv.AudioStatus) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(b.SerialNumber())

	if err := b.MQTTPublishAsync(ctx, MQTTPublishTopicStateMute.Format(sn), s.Mute); err != nil {
		return err
	}

	return b.MQTTPublishAsync(ctx, MQTTPublishTopicStateVolume.Format(sn), s.Volume)
}

func (b *Bind) monitorTvCurrentChannel(s webostv.TvCurrentChannel) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(b.SerialNumber())

	return b.MQTTPublishAsync(ctx, MQTTPublishTopicStateChannelNumber.Format(sn), s.ChannelNumber)
}
