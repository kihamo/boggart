package hikvision

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *ISAPI) StreamingPicture(channel uint64) ([]byte, error) {
	u := fmt.Sprintf("%s/Streaming/channels/%d/picture", a.address, channel)

	request, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	response, err := a.do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return ioutil.ReadAll(response.Body)
}
