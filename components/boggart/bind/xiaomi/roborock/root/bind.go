package root

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
)

var reRuntimeConfigLine = regexp.MustCompile(`(?m)\s*([[:alnum:]_]+)\s*=\s*([^;]+);`)

type Bind struct {
	boggart.BindBase
	di.MQTTBind
	config                 *Config
	cacheRuntimeConfig     map[string]string
	cacheRuntimeConfigLock sync.Mutex
	watchFiles             map[string]func(string) error
}

func (b *Bind) Run() error {
	if b.config.DeviceIDFile != "" {
		if err := b.InitDeviceID(b.config.DeviceIDFile); err != nil {
			return err
		}
	}

	if b.config.RuntimeConfigFile != "" {
		if err := b.AddWatchRuntimeConfig(b.config.RuntimeConfigFile); err != nil {
			return err
		}
	}

	if err := b.StartWatch(); err != nil {
		return err
	}

	for fileName, callback := range b.watchFiles {
		go func(file string, cb func(string) error) {
			if err := cb(file); err != nil {
				b.Logger().Error("Callback returns error",
					"file", file,
					"error", err.Error(),
				)
			}
		}(fileName, callback)
	}

	return nil
}

func (b *Bind) InitDeviceID(fileName string) error {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	b.SetSerialNumber(strings.TrimSpace(string(content)))

	return nil
}

func (b *Bind) AddWatchRuntimeConfig(fileName string) error {
	if _, err := os.Stat(fileName); err != nil {
		return err
	}

	b.watchFiles[fileName] = b.runtimeConfigWatcher

	return nil
}

func (b *Bind) StartWatch() error {
	if len(b.watchFiles) == 0 {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	for file := range b.watchFiles {
		if err = watcher.Add(file); err != nil {
			return err
		}
	}

	go func() {
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
					if w, ok := b.watchFiles[event.Name]; ok {
						if err := w(event.Name); err != nil {
							b.Logger().Error("Watcher callback return error", "error", err.Error(), "file", event.Name)
						}
					}
				}

			case <-watcher.Errors:
				b.Logger().Error("File watcher return error", "error", err.Error())
			}

			// TODO: shutdown
		}
	}()

	return nil
}
