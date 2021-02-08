package nut

import (
	"context"
	"errors"
	"strings"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	variables, err := b.Variables()
	if err != nil {
		return nil, err
	}

	commands, err := b.Commands()
	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	channels := make([]*openhab.Channel, 0, len(variables)+len(commands))

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

	for _, cmd := range commands {
		id := strings.ReplaceAll(cmd.Name, ".", "_")
		id = openhab.IDNormalizeCamelCase("Command " + id)

		channels = append(channels,
			openhab.NewChannel(id, openhab.ChannelTypeSwitch).
				WithStateTopic(b.config.TopicCommand.Format(sn, cmd.Name)).
				WithCommandTopic(b.config.TopicCommandRun.Format(sn, cmd.Name)).
				WithOn(cmd.Name).
				WithOff("done").
				AddItems(
					openhab.NewItem(itemPrefix+id, openhab.ItemTypeSwitch).
						WithLabel(cmd.Description+" []"),
				),
		)
	}

	return openhab.StepsByBind(b, nil, channels...)
}
