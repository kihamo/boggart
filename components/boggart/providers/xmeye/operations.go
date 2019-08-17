package xmeye

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

func (c *Client) OPTime() (*time.Time, error) {
	response := &internal.OPTimeQuery{}

	err := c.CmdWithResult(CmdTimeRequest, "OPTimeQuery", response)
	if err != nil {
		return nil, err
	}

	t, err := time.Parse(timeLayout, response.OPTimeQuery)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *Client) LogSearch(begin time.Time, end time.Time) (interface{}, error) {
	var result interface{}

	err := c.CallWithResult(CmdLogSearchRequest, map[string]interface{}{
		"Name":      "OPLogQuery",
		"SessionID": c.sessionIDAsString(),
		"OPLogQuery": map[string]interface{}{
			"BeginTime":   begin.Format(timeLayout),
			"EndTime":     end.Format(timeLayout),
			"LogPosition": 0,
			"Type":        "LogAll",
		},
	}, &result)
	if err != nil {
		return nil, err
	}

	return result, err
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
