package sun

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() []generators.Step {
	opts, err := b.MQTT().ClientOptions()
	if err != nil {
		return nil
	}

	meta := b.Meta()
	filePrefix := openhab.FilePrefixFromBindMeta(meta)
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)

	broker := openhab.BrokerFromClientOptionsReader(opts)

	steps := []generators.Step{
		{
			FilePath: openhab.DirectoryThings + "broker.things",
			Content:  broker.String(),
		},
	}

	thing := openhab.GenericThingFromBindMeta(meta).
		WithBroker(broker).
		AddChannels(
			openhab.NewChannel("Nadir", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNadir).
				AddItems(
					openhab.NewItem(itemPrefix+"Nadir", openhab.ItemTypeDateTime).
						WithLabel("Nadir").
						WithIcon("time"),
				),
			openhab.NewChannel("NightBeforeStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNightBeforeStart).
				AddItems(
					openhab.NewItem(itemPrefix+"NightBeforeStart", openhab.ItemTypeDateTime).
						WithLabel("Night before start").
						WithIcon("time"),
				),
			openhab.NewChannel("NightBeforeEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNightBeforeEnd).
				AddItems(
					openhab.NewItem(itemPrefix+"NightBeforeEnd", openhab.ItemTypeDateTime).
						WithLabel("Night before end").
						WithIcon("time"),
				),
			openhab.NewChannel("NightBeforeDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicNightBeforeDuration).
				AddItems(
					openhab.NewItem(itemPrefix+"NightBeforeDuration", openhab.ItemTypeNumber).
						WithLabel("Night before duration [%d s]").
						WithIcon("time"),
				),
			openhab.NewChannel("AstronomicalDawnStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicAstronomicalDawnStart),
			openhab.NewChannel("AstronomicalDawnEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicAstronomicalDawnEnd),
			openhab.NewChannel("AstronomicalDawnDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicAstronomicalDawnDuration),
			openhab.NewChannel("DawnStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNauticalDawnStart),
			openhab.NewChannel("DawnEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNauticalDawnEnd),
			openhab.NewChannel("DawnDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicNauticalDawnDuration),
			openhab.NewChannel("CivilDawnStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicCivilDawnStart),
			openhab.NewChannel("CivilDawnEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicCivilDawnEnd),
			openhab.NewChannel("CivilDawnDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicCivilDawnDuration),
			openhab.NewChannel("RiseStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicRiseStart),
			openhab.NewChannel("RiseEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicRiseEnd),
			openhab.NewChannel("RiseDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicRiseDuration),
			openhab.NewChannel("SolarNoon", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicSolarNoon),
			openhab.NewChannel("SetStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicSetStart),
			openhab.NewChannel("SetEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicSetEnd),
			openhab.NewChannel("SetDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicSetDuration),
			openhab.NewChannel("CivilDuskStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicCivilDuskStart),
			openhab.NewChannel("CivilDuskEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicCivilDuskEnd),
			openhab.NewChannel("CivilDuskDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicCivilDuskDuration),
			openhab.NewChannel("NauticalDuskStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNauticalDuskStart),
			openhab.NewChannel("NauticalDuskEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNauticalDuskEnd),
			openhab.NewChannel("NauticalDuskDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicNauticalDuskDuration),
			openhab.NewChannel("AstronomicalDuskStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicAstronomicalDuskStart),
			openhab.NewChannel("AstronomicalDuskEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicAstronomicalDuskEnd),
			openhab.NewChannel("AstronomicalDuskDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicAstronomicalDuskDuration),
			openhab.NewChannel("NightAfterStart", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNightAfterStart),
			openhab.NewChannel("NightAfterEnd", openhab.ChannelTypeDateTime).
				WithStateTopic(b.config.TopicNightAfterEnd),
			openhab.NewChannel("NightAfterDuration", openhab.ChannelTypeNumber).
				WithStateTopic(b.config.TopicNightAfterDuration),
		)

	if content := thing.String(); content != "" {
		steps = append(steps, generators.Step{
			FilePath: openhab.DirectoryThings + filePrefix + ".things",
			Content:  content,
		})
	}

	if content := thing.Items().String(); content != "" {
		steps = append(steps, generators.Step{
			FilePath: openhab.DirectoryItems + filePrefix + ".items",
			Content:  content,
		})
		steps = append(steps, generators.Step{
			FilePath: openhab.DirectoryItems + filePrefix + ".items",
			Content:  content,
		})
	}

	return steps
}
