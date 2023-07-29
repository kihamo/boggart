package timelapse

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	protocol "github.com/kihamo/boggart/protocols/http"
)

const (
	SubDirectoryNameLayout = "2006-01-02"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	di.ProbesBind
	di.WidgetBind
	di.WorkersBind
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	b.Meta().SetLink(&cfg.CaptureURL.URL)

	if cfg.SaveDirectory == "" {
		cacheDir, _ := os.UserCacheDir()
		if cacheDir == "" {
			cacheDir = os.TempDir()
		}

		if cacheDir != "" {
			cacheDirBind := cacheDir + string(os.PathSeparator) + boggart.ComponentName + "_timelapse"

			err := os.Mkdir(cacheDirBind, cfg.DirectoryPerm.FileMode)

			if err == nil {
				b.Logger().Info("Cache dir created", "path", cacheDirBind)
			}

			if err == nil || os.IsExist(err) {
				cacheDir = cacheDirBind
			}
		}

		cfg.SaveDirectory = cacheDir
	}

	return nil
}

func (b *Bind) Capture(ctx context.Context, writer io.Writer) error {
	cfg := b.config()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, cfg.CaptureURL.String(), nil)
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

	subDirectory := filepath.Join(b.config().SaveDirectory, time.Now().Format(SubDirectoryNameLayout))
	fileName := filepath.Join(subDirectory, time.Now().Format(cfg.FileNameFormat)+ext)

	// create sub directory
	if err = os.Mkdir(subDirectory, cfg.DirectoryPerm.FileMode); err != nil && !os.IsExist(err) {
		return err
	}

	fd, err := os.OpenFile(fileName, os.O_CREATE|os.O_EXCL|os.O_WRONLY, cfg.FilePerm.FileMode)
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

func (b *Bind) Files(from, to *time.Time) ([]os.FileInfo, error) {
	saveDirectory := b.config().SaveDirectory

	entries, err := os.ReadDir(saveDirectory)
	if err != nil {
		return nil, err
	}

	var result = make([]os.FileInfo, 0)

	// reverse sorting
	for i := len(entries) - 1; i >= 0; i-- {
		if !entries[i].IsDir() {
			continue
		}

		// filter by directory name
		if from != nil || to != nil {
			ctime, err := time.Parse(SubDirectoryNameLayout, entries[i].Name())
			if err != nil {
				continue
			}

			if from != nil && from.After(time.Date(ctime.Year(), ctime.Month(), ctime.Day(), from.Hour(), from.Minute(), from.Second(), from.Nanosecond(), from.Location())) {
				continue
			}

			if to != nil && to.Before(time.Date(ctime.Year(), ctime.Month(), ctime.Day(), to.Hour(), to.Minute(), to.Second(), to.Nanosecond(), to.Location())) {
				continue
			}
		}

		dir, err := os.Open(filepath.Join(saveDirectory, entries[i].Name()))
		if err != nil {
			return nil, err
		}

		files, err := dir.Readdir(-1)
		if err != nil {
			dir.Close()
			return nil, err
		}

		for _, file := range files {
			// simple filter
			if file.IsDir() || file.Size() == 0 {
				continue
			}

			// by change time
			if from != nil && from.After(file.ModTime()) {
				continue
			}

			if to != nil && to.Before(file.ModTime()) {
				continue
			}

			result = append(result, file)
		}

		dir.Close()
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Name() > result[j].Name() })

	return result, err
}

func (b *Bind) Load(filename string, writer io.Writer) error {
	filename = filepath.Join(b.config().SaveDirectory, filename)
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(writer, f)

	return err
}
