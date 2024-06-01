package mosoblgaz

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

func (b *Bind) InstallerSteps(ctx context.Context, _ installer.System) ([]installer.Step, error) {
	balances, err := b.provider.Balance(ctx)

	if err != nil {
		return nil, err
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	cfg := b.config()

	const (
		idBalance = "Balance"
	)

	channels := make([]*openhab.Channel, 0, len(balances))
	for account, _ := range balances {
		channels = append(channels, openhab.NewChannel(idBalance+account, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicBalance.Format(account)).
			AddItems(
				openhab.NewItem(itemPrefix+idBalance+account, openhab.ItemTypeNumber).
					WithLabel("Balance "+account+" [%.2f ₽]").
					WithIcon("price"),
			))
	}

	return openhab.StepsByBind(b, nil, channels...)
}
