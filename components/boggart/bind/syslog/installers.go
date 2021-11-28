package syslog

import (
	"context"
	"sort"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
	"github.com/kihamo/boggart/components/mqtt"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	publishes := b.MQTT().Publishes()
	channels := make([]*openhab.Channel, 0, len(publishes))

	sortNames := make([]string, 0, len(publishes))
	topics := make(map[string]mqtt.Topic, len(publishes))

	parts := cfg.TopicMessage.Split()

	var nameIndex int

	for i := len(parts) - 1; i > 0; i-- {
		nameIndex = i

		if parts[i] == `+` {
			break
		}
	}

	for topic := range publishes {
		parts := topic.Split()
		if len(parts) <= nameIndex {
			continue
		}
		name := parts[nameIndex]

		if cfg.TopicMessage.Format(name) == topic {
			sortNames = append(sortNames, name)
			topics[name] = topic
		}
	}

	sort.Strings(sortNames)

	for _, name := range sortNames {
		id := openhab.IDNormalizeCamelCase(name)

		channels = append(channels, openhab.NewChannel(id, openhab.ChannelTypeString).
			WithStateTopic(topics[name]).
			AddItems(
				openhab.NewItem(itemPrefix+id, openhab.ItemTypeString).
					WithLabel("Host "+name).
					WithIcon("text"),
			))
	}

	return openhab.StepsByBind(b, nil, channels...)
}
