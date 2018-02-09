package hikvision

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *ISAPI) StreamingPicture(ctx context.Context, channel uint64) ([]byte, error) {
	u := fmt.Sprintf("%s/Streaming/channels/%d/picture", a.address, channel)

	request, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	request = request.WithContext(ctx)

	response, err := a.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
