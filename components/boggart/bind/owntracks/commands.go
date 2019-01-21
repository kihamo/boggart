package owntracks

import (
	"context"
	"encoding/json"
)

var (
	// TODO: setWaypoints, setConfiguration

	commandReportLocation = &Command{
		Type:   "cmd",
		Action: "reportLocation",
	}
	commandRestart = &Command{
		Type:   "cmd",
		Action: "restart",
	}
	commandReconnect = &Command{
		Type:   "cmd",
		Action: "reconnect",
	}
	commandWayPoints = &Command{
		Type:   "cmd",
		Action: "waypoints",
	}
)

func (b *Bind) Command(cmd *Command) error {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return b.MQTTPublish(context.Background(), MQTTOwnTracksPublishTopicCommand.Format(b.user, b.device), 2, false, payload)
}

func (b *Bind) CommandReportLocation() error {
	return b.Command(commandReportLocation)
}

func (b *Bind) CommandRestart() error {
	return b.Command(commandRestart)
}

func (b *Bind) CommandReconnect() error {
	return b.Command(commandReconnect)
}

func (b *Bind) CommandWayPoints() error {
	return b.Command(commandWayPoints)
}
