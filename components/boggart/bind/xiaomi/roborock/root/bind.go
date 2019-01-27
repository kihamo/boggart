package root

import (
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/kihamo/boggart/components/boggart"
)

var reRuntimeConfigLine = regexp.MustCompile(`(?m)\s*([[:alnum:]_]+)\s*=\s*([^;]+);`)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	cacheRuntimeConfig     map[string]string
	cacheRuntimeConfigLock sync.Mutex

	watchFiles map[string]func(string) error
}

func (b *Bind) SetStatusManager(getter boggart.BindStatusGetter, setter boggart.BindStatusSetter) {
	b.BindBase.SetStatusManager(getter, setter)

	b.UpdateStatus(boggart.BindStatusOnline)
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
							// TODO:
							// c.logger.Error("Watcher callback return error", "error", err.Error(), "file", event.Name, )
						}
					}
				}

			case _ = <-watcher.Errors:
				// TODO:
				// c.logger.Error("File watcher return error", "error", err.Error())
			}

			// TODO: shutdown
		}
	}()

	return nil
}
