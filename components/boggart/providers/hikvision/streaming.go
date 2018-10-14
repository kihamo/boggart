package hikvision

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/kihamo/shadow/components/tracing"
	"github.com/opentracing/opentracing-go"
)

func (a *ISAPI) StreamingPicture(ctx context.Context, channel uint64) ([]byte, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, ComponentName+".streaming.picture")
	defer span.Finish()

	if channel < 101 {
		err := fmt.Errorf("Unknown channel %d", channel)

		tracing.SpanError(span, err)
		return nil, err
	}

	u := a.address + "/Streaming/channels/" + strconv.FormatUint(channel, 10) + "/picture"

	response, err := a.Do(ctx, http.MethodGet, u, nil)
	if err != nil {
		tracing.SpanError(span, err)
		return nil, err
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return body, err
}
