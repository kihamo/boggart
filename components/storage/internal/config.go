package internal

import (
	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(storage.ConfigFileNameSpaces, config.ValueTypeString).
			WithUsage("Namespaces").
			WithGroup("File").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add namespace"}),
	}
}
