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

	b.MQTTPublishAsync(ctx, MQTTTopicStateApplication.Format(sn), 2, true, s.AppId)

	// TODO: cache
	if s.AppId == "" {
		b.UpdateStatus(boggart.DeviceStatusOffline)

		b.MQTTPublishAsync(ctx, MQTTTopicStatePower.Format(sn), 2, true, false)
	} else {
		b.MQTTPublishAsync(ctx, MQTTTopicStatePower.Format(sn), 2, true, true)
	}

	return nil
}

func (b *Bind) monitorAudio(s webostv.AudioStatus) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(b.SerialNumber())

	b.MQTTPublishAsync(ctx, MQTTTopicStateMute.Format(sn), 2, true, s.Mute)
	b.MQTTPublishAsync(ctx, MQTTTopicStateVolume.Format(sn), 2, true, s.Volume)

	return nil
}

func (b *Bind) monitorTvCurrentChannel(s webostv.TvCurrentChannel) error {
	ctx := context.Background()
	sn := mqtt.NameReplace(b.SerialNumber())

	b.MQTTPublishAsync(ctx, MQTTTopicStateChannelNumber.Format(sn), 2, true, s.ChannelNumber)

	return nil
}
