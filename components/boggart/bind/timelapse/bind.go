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

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	protocol "github.com/kihamo/boggart/protocols/http"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.ProbesBind
	di.WidgetBind

	config *Config
}

func (b *Bind) Run() error {
	if b.config.SaveDirectory == "" {
		cacheDir, _ := os.UserCacheDir()
		if cacheDir == "" {
			cacheDir = os.TempDir()
		}

		if cacheDir != "" {
			cacheDirBind := cacheDir + string(os.PathSeparator) + boggart.ComponentName + "_timelapse"

			err := os.Mkdir(cacheDirBind, os.FileMode(b.config.SaveDirectoryMode))

			if err == nil {
				b.Logger().Info("Cache dir created", "path", cacheDirBind)
			}

			if err == nil || os.IsExist(err) {
				cacheDir = cacheDirBind
			}
		}

		b.config.SaveDirectory = cacheDir
	}

	return nil
}

func (b *Bind) Capture(ctx context.Context, writer io.Writer) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, b.config.CaptureURL.String(), nil)
	if err != nil {
		return err
	}

	request.Header.Set("User-Agent", protocol.DefaultUserAgent)

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

	fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, os.FileMode(b.config.FileMode))
	if err != nil {
		return err
	}
	defer fd.Close()

	var w io.Writer

	if writer != nil {
		w = io.MultiWriter(fd, writer)
	} else {
		w = fd
	}

	_, err = io.Copy(w, response.Body)

	return err
}
