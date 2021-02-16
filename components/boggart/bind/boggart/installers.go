package boggart

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	return openhab.StepsByBind(b, nil,
		openhab.NewChannel("Name", openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicName.Format(cfg.ApplicationName)).
			AddItems(
				openhab.NewItem(itemPrefix+"Name", openhab.ItemTypeString).
					WithLabel("Name"),
			),
		openhab.NewChannel("Version", openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicBuild.Format(cfg.ApplicationName)).
			AddItems(
				openhab.NewItem(itemPrefix+"Version", openhab.ItemTypeString).
					WithLabel("Version"),
			),
		openhab.NewChannel("Build", openhab.ChannelTypeString).
			WithStateTopic(cfg.TopicVersion.Format(cfg.ApplicationName)).
			AddItems(
				openhab.NewItem(itemPrefix+"Build", openhab.ItemTypeString).
					WithLabel("Build"),
			),
		openhab.NewChannel("Shutdown", openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicShutdown.Format(cfg.ApplicationName)).
			WithCommandTopic(cfg.TopicShutdown.Format(cfg.ApplicationName)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+"Shutdown", openhab.ItemTypeSwitch).
					WithLabel("Shutdown []"),
			),
	)
}
