package xmeye

import (
	"context"
	"io"
	"time"
)

func (c *Client) OPTime(ctx context.Context) (*time.Time, error) {
	var result struct {
		Response
		OPTimeQuery string
	}

	err := c.CmdWithResult(ctx, CmdTimeRequest, "OPTimeQuery", &result)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timeLayout, result.OPTimeQuery)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Client) LogExport(ctx context.Context) (io.Reader, error) {
	packet, err := c.Call(ctx, CmdLogExportRequest, nil)

	return packet.Payload, err
}

// FIXME: после reboot через ручку странное поведение, девайс не перезагружается
// команды принимает, но не отвечает на них
func (c *Client) Reboot(ctx context.Context) error {
	_, err := c.Call(ctx, CmdSysManagerResponse, map[string]interface{}{
		"Name":      "OPMachine",
		"SessionID": c.sessionIDAsString(),
		"OPMachine": map[string]string{
			"Action": "Reboot",
		},
	})

	return err
}
