package internal

import (
	"github.com/kardianos/osext"
	"github.com/kihamo/boggart/components/ota"
	"github.com/kihamo/boggart/components/ota/release"
	"github.com/kihamo/boggart/components/ota/repository"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/dashboard"
)

type Component struct {
	routes []dashboard.Route

	updater          *ota.Updater
	uploadRepository *repository.MemoryRepository
}

func (c *Component) Name() string {
	return ota.ComponentName
}

func (c *Component) Version() string {
	return ota.ComponentVersion
}

func (c *Component) Init(a shadow.Application) error {
	releasePath, err := osext.Executable()
	if err != nil {
		return err
	}

	r, err := release.NewLocalFile(releasePath, a.Version()+" "+a.Build())
	if err != nil {
		return err
	}

	c.updater = ota.NewUpdater()
	c.uploadRepository = repository.NewMemoryRepository()

	c.uploadRepository.Add(r)
	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	return nil
}

func (c *Component) doAutoUpgrade() {
	// TODO:
}
