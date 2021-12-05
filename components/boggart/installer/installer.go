package installer

import (
	"context"
	"path/filepath"
)

type System string
type Path string

const (
	SystemOpenHab System = "OpenHab"
	SystemCron    System = "Cron"
	SystemDevice  System = "Device"
)

func (p Path) Base() string {
	return filepath.Base(string(p))
}

func (p Path) String() string {
	return string(p)
}

type Step struct {
	Description string
	FilePath    Path
	Content     string
}

type HasInstaller interface {
	InstallersSupport() []System
	InstallerSteps(context.Context, System) ([]Step, error)
}
