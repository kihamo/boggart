package telegram

import (
	"context"
	"strconv"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/go-workers"
	"github.com/kihamo/go-workers/task"
	"gopkg.in/telegram-bot-api.v4"
)

func (b *Bind) Tasks() []workers.Task {
	taskLiveness := task.NewFunctionTask(b.taskLiveness)
	taskLiveness.SetTimeout(b.config.LivenessTimeout)
	taskLiveness.SetRepeats(-1)
	taskLiveness.SetRepeatInterval(b.config.LivenessInterval)
	taskLiveness.SetName("liveness-" + b.config.Token)

	return []workers.Task{
		taskLiveness,
	}
}

func (b *Bind) taskLiveness(ctx context.Context) (interface{}, error) {
	var needInit bool

	bot := b.bot()
	if bot != nil {
		if _, err := bot.GetMe(); err != nil {
			needInit = true

			b.mutex.Lock()
			b.client = nil
			b.mutex.Unlock()

			b.UpdateStatus(boggart.BindStatusOffline)

			bot.StopReceivingUpdates()
			b.done <- struct{}{}
		}
	} else {
		needInit = true

		b.UpdateStatus(boggart.BindStatusOffline)
	}

	if !needInit {
		b.UpdateStatus(boggart.BindStatusOnline)
		return nil, nil
	}

	client, err := tgbotapi.NewBotAPI(b.config.Token)
	if err != nil {
		return nil, err
	}

	b.SetSerialNumber(strconv.Itoa(client.Self.ID))
	client.Debug = b.config.Debug

	if b.config.UpdatesEnabled {
		client.Buffer = b.config.UpdatesBuffer

		u := tgbotapi.NewUpdate(0)
		u.Timeout = b.config.UpdatesTimeout

		updates, err := client.GetUpdatesChan(u)
		if err != nil {
			return nil, err
		}

		b.listenUpdates(updates)
	}

	b.mutex.Lock()
	b.client = client
	b.mutex.Unlock()

	b.UpdateStatus(boggart.BindStatusOnline)

	return nil, nil
}
