package probes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/kihamo/boggart/components/boggart"
)

func HTTPProbe(ctx context.Context, address string, headers http.Header) error {
	u, err := url.Parse(address)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	if _, ok := headers["User-Agent"]; !ok {
		if headers == nil {
			headers = http.Header{}
		}
		headers.Set("User-Agent", boggart.ComponentName+"/"+boggart.ComponentVersion)
	}

	req.Header = headers
	if headers.Get("Host") != "" {
		req.Host = headers.Get("Host")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusOK && res.StatusCode < http.StatusBadRequest {
		if res.StatusCode >= http.StatusMultipleChoices {
			return fmt.Errorf("HTTP probe failed with statuscode: %d", res.StatusCode)
		}

		return nil
	}

	return fmt.Errorf("HTTP probe failed with statuscode: %d", res.StatusCode)
}
