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

	t, err := time.Parse(TimeLayout, result.OPTimeQuery)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Client) OPTimeSetting(ctx context.Context, t time.Time) error {
	_, err := c.Call(ctx, CmdSysManagerRequest, map[string]interface{}{
		"Name":          "OPTimeSetting",
		"SessionID":     c.connection.SessionIDAsString(),
		"OPTimeSetting": t.Format(TimeLayout),
	})

	return err
}

func (c *Client) LogExport(ctx context.Context) (io.Reader, error) {
	packet, err := c.Call(ctx, CmdLogExportRequest, nil)

	return packet.payload, err
}

// FIXME: после reboot через ручку странное поведение, девайс не перезагружается
// команды принимает, но не отвечает на них
func (c *Client) Reboot(ctx context.Context) error {
	_, err := c.Call(ctx, CmdSysManagerResponse, map[string]interface{}{
		"Name":      "OPMachine",
		"SessionID": c.connection.SessionIDAsString(),
		"OPMachine": map[string]string{
			"Action": "Reboot",
		},
	})

	return err
}
