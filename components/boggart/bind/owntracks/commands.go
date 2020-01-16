package owntracks

import (
	"context"
	"encoding/json"
)

func (b *Bind) sendCommand(cmd interface{}) error {
	payload, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return b.MQTTContainer().Publish(context.Background(), b.config.TopicOwnTracksCommand, payload)
}

func (b *Bind) CommandReportLocation() error {
	return b.sendCommand(&CommandPayload{
		Type:   "cmd",
		Action: "reportLocation",
	})
}

func (b *Bind) CommandRestart() error {
	return b.sendCommand(&CommandPayload{
		Type:   "cmd",
		Action: "restart",
	})
}

func (b *Bind) CommandReconnect() error {
	return b.sendCommand(&CommandPayload{
		Type:   "cmd",
		Action: "reconnect",
	})
}

func (b *Bind) CommandWayPoints() error {
	return b.sendCommand(&CommandPayload{
		Type:   "cmd",
		Action: "waypoints",
	})
}

func (b *Bind) CommandSetWayPoints(points []WayPointPayload) error {
	for i := range points {
		points[i].Type = "waypoint"
	}

	return b.sendCommand(&SetWayPointsPayload{
		CommandPayload: CommandPayload{
			Type:   "cmd",
			Action: "setWaypoints",
		},
		WayPoints: WayPointsPayload{
			Type:      "waypoints",
			WayPoints: points,
		},
	})
}
