package timelapse

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"sort"
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

			err := os.Mkdir(cacheDirBind, b.config.SaveDirectoryMode.FileMode)

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

	fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, b.config.FileMode.FileMode)
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

func (b *Bind) Files() ([]os.FileInfo, error) {
	dir, err := os.Open(b.config.SaveDirectory)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err == nil {
		for i := len(files) - 1; i >= 0; i-- {
			if files[i].IsDir() || files[i].Size() == 0 {
				files = append(files[:i], files[i+1:]...)
			}
		}
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Name() > files[j].Name() })

	return files, err
}

func (b *Bind) Load(filename string, writer io.Writer) error {
	filename = b.config.SaveDirectory + string(os.PathSeparator) + filename
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(writer, f)

	return err
}
