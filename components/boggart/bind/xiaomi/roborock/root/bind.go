package root

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/performance"
)

var reRuntimeConfigLine = regexp.MustCompile(`(?m)\s*([[:alnum:]_]+)\s*=\s*([^;]+);`)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MQTTBind

	cacheRuntimeConfig     map[string]string
	cacheRuntimeConfigLock sync.Mutex
	watchFiles             map[string]func(string) error
}

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	b.cacheRuntimeConfig = make(map[string]string, 11)
	b.watchFiles = make(map[string]func(string) error)

	cfg := b.config()

	if cfg.DeviceIDFile != "" {
		if err := b.InitDeviceID(cfg.DeviceIDFile); err != nil {
			return err
		}
	}

	if cfg.RuntimeConfigFile != "" {
		if err := b.AddWatchRuntimeConfig(cfg.RuntimeConfigFile); err != nil {
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

	b.Meta().SetSerialNumber(performance.UnsafeBytes2String(bytes.TrimSpace(content)))

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
		}
	}()

	return nil
}
