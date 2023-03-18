package xmeye

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
)

func (c *Client) Snapshot(ctx context.Context, channel int) ([]byte, error) {
	packet, err := c.Call(ctx, CmdSnapshot, map[string]interface{}{
		"Name":      "OPSNAP",
		"SessionID": c.connection.SessionIDAsString(),
		"OPSNAP": map[string]interface{}{
			"Channel": channel,
		},
	})

	if err != nil {
		return nil, err
	}

	return packet.payload.Bytes(), nil
}

func (c *Client) SnapshotImage(ctx context.Context, channel int) (image.Image, error) {
	b, err := c.Snapshot(ctx, channel)
	if err != nil {
		return nil, err
	}
	return jpeg.Decode(bytes.NewReader(b))
}
