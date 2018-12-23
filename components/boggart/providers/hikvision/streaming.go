package hikvision

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/kihamo/shadow/components/tracing"
)

func (a *ISAPI) StreamingPictureToWriter(ctx context.Context, channel uint64, writer io.Writer) error {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "streaming.picture")
	defer span.Finish()

	if channel < 101 {
		err := fmt.Errorf("unknown channel %d", channel)

		tracing.SpanError(span, err)
		return err
	}

	u := a.address + "/Streaming/channels/" + strconv.FormatUint(channel, 10) + "/picture"

	response, err := a.Do(ctx, http.MethodGet, u, nil)
	if err != nil {
		tracing.SpanError(span, err)
		return err
	}

	_, err = io.Copy(writer, response.Body)
	response.Body.Close()

	if err != nil {
		tracing.SpanError(span, err)
	}

	return err
}

func (a *ISAPI) StreamingPicture(ctx context.Context, channel uint64) ([]byte, error) {
	buf := &bytes.Buffer{}

	err := a.StreamingPictureToWriter(ctx, channel, buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
