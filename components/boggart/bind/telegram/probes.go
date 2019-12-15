package telegram

import (
	"context"
)

func (b *Bind) ReadinessProbe(ctx context.Context) (err error) {
	bot := b.bot()

	if bot == nil {
		bot, err = b.initBot()
		if err != nil {
			return err
		}
	}

	_, err = bot.GetMe()

	return err
}

func (b *Bind) LivenessProbe(ctx context.Context) (err error) {
	b.mutex.RLock()
	client := b.client
	b.mutex.RUnlock()

	if client == nil {
		return nil
	}

	_, err = client.GetMe()

	return err
}
