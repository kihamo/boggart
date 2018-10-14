package hikvision

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	tracing "github.com/kihamo/shadow/components/tracing/http"
)

func (a *ISAPI) StreamingPicture(ctx context.Context, channel uint64) ([]byte, error) {
	if channel < 101 {
		return nil, fmt.Errorf("Unknown channel %d", channel)
	}

	u := a.address + "/Streaming/channels/" + strconv.FormatUint(channel, 10) + "/picture"

	ctx = tracing.OperationNameToContext(ctx, ComponentName+".StreamingPicture")

	response, err := a.Do(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
