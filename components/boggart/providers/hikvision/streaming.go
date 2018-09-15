package hikvision

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

func (a *ISAPI) StreamingPicture(ctx context.Context, channel uint64) ([]byte, error) {
	if channel < 101 {
		return nil, fmt.Errorf("Unknown channel %d", channel)
	}

	u := a.address + "/Streaming/channels/" + strconv.FormatUint(channel, 10) + "/picture"

	response, err := a.Do(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
