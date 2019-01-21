package owntracks

import (
	"context"
	"encoding/json"
)

type Command struct {
	Type   string `json:"_type"`
	Action string `json:"action"`
}

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

func (b *Bind) Command(user, device string, cmd *Command) error {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return b.MQTTPublish(context.Background(), MQTTSubscribeTopicCommand.Format(user, device), 2, false, payload)
}

func (b *Bind) CommandReportLocation(user, device string) error {
	return b.Command(user, device, commandReportLocation)
}

func (b *Bind) CommandRestart(user, device string) error {
	return b.Command(user, device, commandRestart)
}

func (b *Bind) CommandReconnect(user, device string) error {
	return b.Command(user, device, commandReconnect)
}

func (b *Bind) CommandWayPoints(user, device string) error {
	return b.Command(user, device, commandWayPoints)
}
