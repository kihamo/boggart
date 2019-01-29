package grafana

import (
	"context"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	g "github.com/kihamo/go-grafana-api"
	"go.uber.org/multierr"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	client     *g.Client
	name       string
	dashboards []int64
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
}

func (b *Bind) CreateAnnotation(title, text string, tags []string, timeStart *time.Time, timeEnd *time.Time) (err error) {
	body := []string{title, text}

	if timeStart == nil {
		now := time.Now()
		timeStart = &now
	}

	input := &g.CreateAnnotationInput{
		Time: g.Int64(timeStart.UnixNano() / int64(time.Millisecond)),
		Text: g.String(strings.Join(body, "\n")),
		Tags: g.StringSlice(tags),
	}

	if timeEnd != nil {
		input.TimeEnd = g.Int64(timeEnd.UnixNano() / int64(time.Millisecond))
		input.IsRegion = g.Bool(true)
	}

	for _, dashboardId := range b.dashboards {
		input.DashboardId = g.Int64(dashboardId)

		if _, e := b.client.CreateAnnotation(context.Background(), input); e != nil {
			err = multierr.Append(err, e)
		}
	}

	return err
}