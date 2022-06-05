package pass24online

import (
	"context"
	"encoding/json"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/providers/pass24online/client/feed"
	"go.uber.org/multierr"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), b.config().UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	response, err := b.provider.Feed.GetFeed(feed.NewGetFeedParams().WithFilterType(&[]uint64{4}[0]), nil)
	if err != nil {
		return err
	}

	cfg := b.config()
	feedStartDatetime := b.feedStartDatetime.Load()
	collection := response.GetPayload().Body.Collection

	for i := len(collection) - 1; i >= 0; i-- {
		dt := collection[i].HappenedAt.Time()
		if dt.Before(feedStartDatetime) {
			continue
		}

		event := FeedEvent{
			ModelName:   collection[i].Subject.GuestData.Model.Name,
			PlateNumber: collection[i].Subject.GuestData.PlateNumber,
			Message:     collection[i].Message,
			Datetime:    dt,
		}

		if raw, ok := collection[i].EventData["status"]; ok {
			if val, ok := raw.(json.Number); ok {
				if name, ok := statusName[val.String()]; ok {
					event.Status = name
				}
			}
		}

		if event.Status == "" {
			event.Status = statusName["0"]
		}

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicFeedEvent.Format(cfg.Phone), event); e != nil {
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
