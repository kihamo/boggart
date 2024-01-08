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
	claim.messageID = CmdPlayClaimRequest
	err = claim.LoadPayload(map[string]interface{}{
		"Name":      "OPPlayBack",
		"SessionID": c.connection.SessionIDAsString(),
		"OPPlayBack": map[string]interface{}{
			"Action":     "Claim",
			"StartTime":  begin.Format(TimeLayout),
			"EndTime":    end.Format(TimeLayout),
			"StreamType": 0,
			"Parameter": map[string]interface{}{
				"FileName":  name,
				"TransMode": "TCP",
				// more fields https://github.com/sahujaunpuri/XMCamera-AndroidTV-Demo/blob/b830496cc59d5e50105dfa99250e5e1cb324b28f/app/src/main/java/com/lib/sdk/bean/OPPlayBackBean.java#L19
				"PlayMode":                 "ByName",
				"Value":                    0,
				"IntelligentPlayBackEvent": "",
				"IntelligentPlayBackSpeed": 0,
			},
		},
	})

	if err != nil {
		return nil, err
	}

	request := newPacket()
	request.messageID = CmdPlayRequest
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

	dial, err := protocol.NewByDSNString(c.dsn)
	if err != nil {
		return nil, err
	}

	return newStream(&connection{
		Connection:     dial,
		sessionID:      c.connection.SessionID(),
		sequenceNumber: 0,
	}, claim, func() error {
		return c.connection.send(request)
	}), nil
}
