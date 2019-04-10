package hikvision

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/kihamo/shadow/components/tracing"
)

const (
	imagePrefixURL = "/Image"

	ImageIrCutFilterTypeAuto         ImageIrCutFilterType = "auto"
	ImageIrCutFilterTypeDay          ImageIrCutFilterType = "day"
	ImageIrCutFilterTypeNight        ImageIrCutFilterType = "night"
	ImageIrCutFilterTypeSchedule     ImageIrCutFilterType = "schedule"
	ImageIrCutFilterTypeEventTrigger ImageIrCutFilterType = "eventTrigger"
)

type ImageIrCutFilterType string

type ImageChannelList struct {
	Channels []ImageChannel `xml:"ImageChannel"`
}

type ImageChannel struct {
	ID             uint64               `xml:"id"`
	Enabled        bool                 `xml:"enabled"`
	VideoInputID   uint64               `xml:"videoInputID"`
	IrCutFilter    *ImageIrCutFilter    `xml:"IrcutFilter"`
	IrCutFilterExt *ImageIrCutFilterExt `xml:"IrcutFilterExt"`
}

type ImageIrCutFilter struct {
	Type  ImageIrCutFilterType `xml:"IrcutFilterType"`
	Level string               `xml:"IrcutFilterLevel"`
	Time  uint64               `xml:"IrcutFilterTime"`
}

type ImageIrCutFilterExt struct {
	Type                  ImageIrCutFilterType      `xml:"IrcutFilterType"`
	NightToDayFilterLevel string                    `xml:"nightToDayFilterLevel"`
	NightToDayFilterTime  uint64                    `xml:"nightToDayFilterTime"`
	Schedule              *ImageIrCutFilterSchedule `xml:"Schedule"`
}

type ImageIrCutFilterSchedule struct {
	Type      string                            `xml:"scheduleType"`
	TimeRange ImageIrCutFilterScheduleTimeRange `xml:"TimeRange"`
}

type ImageIrCutFilterScheduleTimeRange struct {
	Begin string `xml:"beginTime"`
	End   string `xml:"endTime"`
}

func (a *ISAPI) ImageChannels(ctx context.Context) (list ImageChannelList, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "image.channels")
	defer span.Finish()

	u := a.address + imagePrefixURL + "/channels"

	err = a.DoXML(ctx, http.MethodGet, u, nil, &list)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return list, err
}

func (a *ISAPI) ImageIrCutFilter(ctx context.Context, channel uint64, filter ImageIrCutFilter) error {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "image.ir_cut_filter")
	defer span.Finish()

	u := a.address + imagePrefixURL + "/channels/" + strconv.FormatUint(channel, 10) + "/IrcutFilter"

	result := ResponseStatus{}

	err := a.DoXML(ctx, http.MethodPut, u, filter, &result)
	if err != nil {
		tracing.SpanError(span, err)
		return err
	}

	if result.StatusCode != 1 {
		err := errors.New(result.StatusString)

		tracing.SpanError(span, err)
		return err
	}

	return nil
}
