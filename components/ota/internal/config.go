package internal

import (
	"os"

	"github.com/kihamo/boggart/components/ota"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(ota.ConfigReleasesDirectory, config.ValueTypeString).
			WithUsage("Path to saved releases directory").
			WithEditable(true).
			WithDefault(os.TempDir()),
		config.NewVariable(ota.ConfigRepositoryServerEnabled, config.ValueTypeBool).
			WithUsage("Enable serve repository").
			WithEditable(true).
			WithDefault(true),
	}
}
