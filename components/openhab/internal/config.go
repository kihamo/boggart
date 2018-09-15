package internal

import (
	"github.com/kihamo/boggart/components/openhab"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) ConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(openhab.ConfigAPIURL, config.ValueTypeString).
			WithUsage("API URL").
			WithGroup("API").
			WithEditable(true),
		config.NewVariable(openhab.ConfigTelegramChats, config.ValueTypeString).
			WithUsage("Chats for messages").
			WithGroup("Messenger Telegram").
			WithEditable(true).
			WithView([]string{config.ViewTags}).
			WithViewOptions(map[string]interface{}{config.ViewOptionTagsDefaultText: "add a chat ID"}),
	}
}
