package xmeye

import (
	"io"
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

func (c *Client) OPTime() (*time.Time, error) {
	var result struct {
		internal.Response
		OPTimeQuery string
	}

	err := c.CmdWithResult(CmdTimeRequest, "OPTimeQuery", &result)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timeLayout, result.OPTimeQuery)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Client) LogExport() (io.Reader, error) {
	packet, err := c.Call(CmdLogExportRequest, nil)

	return packet.Payload, err
}

// FIXME: после reboot через ручку странное поведение, девайс не перезагружается
// команды принимает, но не отвечает на них
func (c *Client) Reboot() error {
	_, err := c.Call(CmdSysManagerResponse, map[string]interface{}{
		"Name":      "OPMachine",
		"SessionID": c.sessionIDAsString(),
		"OPMachine": map[string]string{
			"Action": "Reboot",
		},
	})

	return err
}
