package internal

import (
	"net/url"

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

func (c *Component) ConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher([]string{
			openhab.ConfigAPIURL,
		}, c.watchAPIURL),
	}
}

func (c *Component) watchAPIURL(_ string, newValue interface{}, _ interface{}) {
	if apiUrl, err := url.Parse(newValue.(string)); err == nil {
		c.mutex.Lock()
		c.apiUrl = apiUrl
		c.mutex.Unlock()
	}
}
