package mosoblgaz

import (
	"context"

	"github.com/hashicorp/go-multierror"
)

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	balances, err := b.provider.Balance(ctx)

	if err != nil {
		return err
	}

	cfg := b.config()

	for account, balance := range balances {
		metricBalance.With("account", account).Set(balance)

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicBalance.Format(account), balance); e != nil {
			err = multierror.Append(err, e)
		}
	}

	return nil
}
