package internal

import (
	"github.com/fsnotify/fsnotify"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/roborock"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/logging"
)

type Component struct {
	application shadow.Application
	logger      logging.Logger
	mqtt        mqtt.Component

	files map[string]func(string) error
}

func (c *Component) Name() string {
	return roborock.ComponentName
}

func (c *Component) Version() string {
	return roborock.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name: logging.ComponentName,
		},
		{
			Name:     mqtt.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a
	c.mqtt = a.GetComponent(mqtt.ComponentName).(mqtt.Component)
	c.files = map[string]func(string) error{
		roborock.FileRuntimeConfig: c.runtimeConfigWatcher,
	}

	return nil
}

func (c *Component) Run() error {
	c.logger = logging.DefaultLogger().Named(c.Name())

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	// watch file
	for file := range c.files {
		if err = watcher.Add(file); err != nil {
			return err
		}
	}

	for {
		select {
		case event := <-watcher.Events:
			if event.Op&fsnotify.Write == fsnotify.Write {
				c.logger.Debug("Watched file modified", "file", event.Name)

				// call watcher
				if w, ok := c.files[event.Name]; ok {
					if err := w(event.Name); err != nil {
						c.logger.Error("Watcher callback return error",
							"error", err.Error(),
							"file", event.Name,
						)
					}
				}
			}

		case err := <-watcher.Errors:
			c.logger.Error("File watcher return error", "error", err.Error())
		}
	}

	return nil
}
