package nut

import (
	"errors"
	"strings"

	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	variables, err := b.Variables()
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	channels := make([]*openhab.Channel, 0, 1+len(variables))

	channels = append(channels,
		openhab.NewChannel("Command", openhab.ChannelTypeString).
			WithStateTopic(b.config.TopicCommand.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+"Command", openhab.ItemTypeString).
					WithLabel("Command"),
			),
	)

	if variables, err := b.Variables(); err == nil {
		var channel *openhab.Channel

		for _, v := range variables {
			id := strings.ReplaceAll(v.Name, ".", "_")
			id = openhab.IDNormalizeCamelCase(id)

			switch v.Value.(type) {
			case float64, int:
				channel = openhab.NewChannel(id, openhab.ChannelTypeNumber).
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeNumber).
							WithLabel(v.Description),
					)
			case bool:
				channel = openhab.NewChannel(id, openhab.ChannelTypeContact).
					WithOn("true").
					WithOff("false").
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeContact).
							WithLabel(v.Description),
					)
			default:
				channel = openhab.NewChannel(id, openhab.ChannelTypeString).
					AddItems(
						openhab.NewItem(itemPrefix+id, openhab.ItemTypeString).
							WithLabel(v.Description),
					)
			}

			channel.WithStateTopic(b.config.TopicVariable.Format(sn, v.Name))
			if v.Type.Writeable {
				channel.WithCommandTopic(b.config.TopicVariableSet.Format(sn, v.Name))
			}

			channels = append(channels, channel)
		}
	}

	return openhab.StepsByBind(b, nil, channels...)
}
