package telegram

import (
	"context"
)

func (b *Bind) LivenessProbe(_ context.Context) (err error) {
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

func (b *Bind) ReadinessProbe(_ context.Context) (err error) {
	if bot := b.bot(); bot != nil {
		_, err = bot.GetMe()
	}

	return err
}
