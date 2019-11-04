package octoprint

import (
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := b.config.Address.Host

	return []mqtt.Topic{
		b.config.TopicState.Format(sn),
		b.config.TopicStateBedTemperatureActual.Format(sn),
		b.config.TopicStateBedTemperatureOffset.Format(sn),
		b.config.TopicStateBedTemperatureTarget.Format(sn),
		b.config.TopicStateTool0TemperatureActual.Format(sn),
		b.config.TopicStateTool0TemperatureOffset.Format(sn),
		b.config.TopicStateTool0TemperatureTarget.Format(sn),
		b.config.TopicStateJobFileName.Format(sn),
		b.config.TopicStateJobFileSize.Format(sn),
		b.config.TopicStateJobProgress.Format(sn),
		b.config.TopicStateJobTime.Format(sn),
		b.config.TopicStateJobTimeLeft.Format(sn),
		b.config.TopicLayerTotal.Format(sn),
		b.config.TopicLayerCurrent.Format(sn),
		b.config.TopicHeightTotal.Format(sn),
		b.config.TopicHeightCurrent.Format(sn),
	}
}
