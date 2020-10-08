package pass24online

import (
	"context"
	"time"

	"github.com/kihamo/boggart/providers/pass24online/client/feed"
	"github.com/kihamo/go-workers"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []workers.Task {
	taskStateUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskStateUpdater.SetRepeats(-1)
	taskStateUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskStateUpdater.SetName("updater")

	return []workers.Task{
		taskStateUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	response, err := b.provider.Feed.GetFeed(feed.NewGetFeedParams().WithFilterType(&[]uint64{4}[0]), nil)
	if err != nil {
		return err
	}

	feedStartDatetime := b.feedStartDatetime.Load()
	collection := response.GetPayload().Body.Collection

	for i := len(collection) - 1; i >= 0; i-- {
		dt := collection[i].HappenedAt.Time()
		if dt.Before(feedStartDatetime) {
			continue
		}

		if e := b.MQTT().PublishAsync(ctx, b.config.TopicFeedEvent.Format(b.config.Phone), collection[i].Message); e != nil {
			err = multierr.Append(err, e)
		}

		b.feedStartDatetime.Set(dt)
		feedStartDatetime = dt
	}

	// так как у нас на равенство проверка, то в следующие итерацию события так же
	// попадут в рассылку если не было новых, для избежания этого прибавляем секунду
	b.feedStartDatetime.Set(feedStartDatetime.Add(time.Second))

	return err
}
