package internal

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/components/roborock"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/logger"
)

type Component struct {
	application shadow.Application
	config      config.Component
	logger      logger.Logger
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
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: logger.ComponentName,
		},
		{
			Name:     mqtt.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) (err error) {
	c.application = a
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	c.mqtt = a.GetComponent(mqtt.ComponentName).(mqtt.Component)
	c.files = map[string]func(string) error{
		roborock.FileRuntimeConfig: c.runtimeConfigWatcher,
	}

	return nil
}

func (c *Component) Run(wg *sync.WaitGroup) error {
	c.logger = logger.NewOrNop(c.Name(), c.application)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer wg.Done()

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					c.logger.Debugf("Watched file modified", map[string]interface{}{
						"file": event.Name,
					})

					// call watcher
					if w, ok := c.files[event.Name]; ok {
						if err := w(event.Name); err != nil {
							c.logger.Error("Watcher callback return error", map[string]interface{}{
								"error": err.Error(),
								"file":  event.Name,
							})
						}
					}
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}

				c.logger.Error("File watcher return error", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}
	}()

	// watch file
	for file := range c.files {
		if err = watcher.Add(file); err != nil {
			return err
		}
	}

	return nil
}
