package timelapse

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart/di"
)

const (
	DefaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/63.0.3239.132 Safari/537.36"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.ProbesBind

	config *Config
}

func (b *Bind) Capture(ctx context.Context) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, b.config.CaptureURL.String(), nil)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", DefaultUserAgent)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("capture failed with statuscode: %d", response.StatusCode)
	}

	var ext string

	if contentType := response.Header.Get("Content-type"); contentType != "" {
		for _, v := range strings.Split(contentType, ",") {
			if m, _, err := mime.ParseMediaType(v); err == nil {
				switch m {
				case "image/gif":
					ext = ".gif"

				case "image/jpeg":
					ext = ".jpg"

				case "image/png":
					ext = ".png"
				}

				break
			}
		}
	}

	fileName := b.config.SaveDirectory + string(os.PathSeparator) + time.Now().Format(b.config.FileNameFormat) + ext

	fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0660)
	if err != nil {
		return err
	}
	defer fd.Close()

	_, err = io.Copy(fd, response.Body)

	return err
}
