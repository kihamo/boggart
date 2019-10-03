package xmeye

import (
	"context"
	"io"
	"time"

	protocol "github.com/kihamo/boggart/protocols/connection"
)

func (c *Client) PlayStream(ctx context.Context, begin, end time.Time, name string) (io.ReadCloser, error) {
	_, err := c.Call(ctx, CmdPlayRequest, map[string]interface{}{
		"Name":      "OPPlayBack",
		"SessionID": c.connection.SessionIDAsString(),
		"OPPlayBack": map[string]interface{}{
			"Action":    "DownloadStart",
			"StartTime": begin.Format(TimeLayout),
			"EndTime":   end.Format(TimeLayout),
			"Parameter": map[string]interface{}{
				"FileName":  name,
				"TransMode": "TCP",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	claim := newPacket()
	claim.MessageID = CmdPlayClaimRequest
	err = claim.LoadPayload(map[string]interface{}{
		"Name":      "OPPlayBack",
		"SessionID": c.connection.SessionIDAsString(),
		"OPPlayBack": map[string]interface{}{
			"Action":    "Claim",
			"StartTime": begin.Format(TimeLayout),
			"EndTime":   end.Format(TimeLayout),
			"Parameter": map[string]interface{}{
				"FileName":  name,
				"TransMode": "TCP",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	request := newPacket()
	request.MessageID = CmdPlayRequest
	err = request.LoadPayload(map[string]interface{}{
		"Name":      "OPPlayBack",
		"SessionID": c.connection.SessionIDAsString(),
		"OPPlayBack": map[string]interface{}{
			"Action":    "DownloadStart",
			"StartTime": begin.Format(TimeLayout),
			"EndTime":   end.Format(TimeLayout),
			"Parameter": map[string]interface{}{
				"FileName":  name,
				"TransMode": "TCP",
			},
		},
	})

	if err != nil {
		return nil, err
	}

	dial, err := protocol.New(c.dsn)
	if err != nil {
		return nil, err
	}

	return newStream(&connection{
		Conn:           dial,
		sessionID:      c.connection.SessionID(),
		sequenceNumber: 0,
	}, claim, func() error {
		return c.connection.send(request)
	}), nil
}
