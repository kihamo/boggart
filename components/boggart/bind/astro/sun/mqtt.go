package sun

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	MQTTPrefix mqtt.Topic = boggart.ComponentName + "/astro/sun/+/"

	// восход солнца (верхний край солнца появляется на горизонте)
	MQTTPublishTopicRiseStart = MQTTPrefix + "rise/start"
	// конец восхода (нижний край солнца касается горизонта)
	MQTTPublishTopicRiseEnd      = MQTTPrefix + "rise/end"
	MQTTPublishTopicRiseDuration = MQTTPrefix + "rise/duration"
	// начинается закат (нижний край солнца касается горизонта)
	MQTTPublishTopicSetStart = MQTTPrefix + "set/start"
	// закат (солнце исчезает за горизонтом, начинается вечерние гражданские сумерки)
	MQTTPublishTopicSetEnd      = MQTTPrefix + "set/end"
	MQTTPublishTopicSetDuration = MQTTPrefix + "set/duration"
	// начинается ночь (достаточно темная для астрономических наблюдений)
	MQTTPublishTopicNightStart = MQTTPrefix + "night/start"
	// кончается ночь (начинаются утренние астрономические сумерки)
	MQTTPublishTopicNightEnd      = MQTTPrefix + "night/end"
	MQTTPublishTopicNightDuration = MQTTPrefix + "night/duration"
	// надир (самый темный момент ночи, солнце находится в самой низкой позиции)
	MQTTPublishTopicNadir = MQTTPrefix + "nadir"
	// солнечный полдень (солнце в самом верхнем положении)
	MQTTPublishTopicSolarNoon = MQTTPrefix + "solar-noon"
)

func (b *Bind) MQTTPublishes() []mqtt.Topic {
	sn := mqtt.NameReplace(b.SerialNumber())

	return []mqtt.Topic{
		mqtt.Topic(MQTTPublishTopicRiseStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicRiseEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicRiseDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSetStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSetEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSetDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNightStart.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNightEnd.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNightDuration.Format(sn)),
		mqtt.Topic(MQTTPublishTopicNadir.Format(sn)),
		mqtt.Topic(MQTTPublishTopicSolarNoon.Format(sn)),
	}
}
