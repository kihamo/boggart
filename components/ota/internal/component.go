package internal

import (
	"github.com/kardianos/osext"
	"github.com/kihamo/boggart/components/ota"
	"github.com/kihamo/boggart/components/ota/release"
	"github.com/kihamo/boggart/components/ota/repository"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
	"github.com/kihamo/shadow/components/logging"
)

type Component struct {
	config config.Component
	routes []dashboard.Route

	updater          *ota.Updater
	uploadRepository *repository.DirectoryRepository
	currentRelease   ota.Release
}

func (c *Component) Name() string {
	return ota.ComponentName
}

func (c *Component) Version() string {
	return ota.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
		{
			Name: dashboard.ComponentName,
		},
		{
			Name: i18n.ComponentName,
		},
		{
			Name: logging.ComponentName,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	releasePath, err := osext.Executable()
	if err != nil {
		return err
	}

	c.currentRelease, err = release.NewLocalFile(releasePath, a.Version()+" "+a.Build())
	if err != nil {
		return err
	}

	c.updater = ota.NewUpdater()

	c.uploadRepository = repository.NewDirectoryRepository()
	c.uploadRepository.Add(c.currentRelease)

	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	<-a.ReadyComponent(config.ComponentName)
	cfg := a.GetComponent(config.ComponentName).(config.Component)

	return c.uploadRepository.Load(cfg.String(ota.ConfigReleasesDirectory))
}

func (c *Component) doAutoUpgrade() {
	// TODO:
}
