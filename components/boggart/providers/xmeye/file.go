package xmeye

import (
	"time"

	"github.com/kihamo/boggart/components/boggart/providers/xmeye/internal"
)

type fileSearchEvent string
type fileSearchType string

const (
	FileSearchEventAll          fileSearchEvent = "*"
	FileSearchEventAlarm        fileSearchEvent = "A"
	FileSearchEventMotionDetect fileSearchEvent = "M"
	FileSearchEventGeneral      fileSearchEvent = "R"
	FileSearchEventManual       fileSearchEvent = "H"

	FileSearchH264 fileSearchType = "h264"
	FileSearchJPEG fileSearchType = "jpg"
)

func (c *Client) LogSearch(begin, end time.Time) ([]internal.LogSearch, error) {
	var result struct {
		internal.Response
		OPLogQuery []internal.LogSearch
	}

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

	return result.OPLogQuery, err
}

func (c *Client) FileSearch(begin, end time.Time, channel uint32, event fileSearchEvent, typ fileSearchType) ([]internal.FileSearch, error) {
	var result struct {
		internal.Response
		OPFileQuery []internal.FileSearch
	}

	err := c.CallWithResult(CmdFileSearchRequest, map[string]interface{}{
		"Name":      "OPFileQuery",
		"SessionID": c.sessionIDAsString(),
		"OPFileQuery": map[string]interface{}{
			"BeginTime":      begin.Format(timeLayout),
			"EndTime":        end.Format(timeLayout),
			"Channel":        channel,
			"DriverTypeMask": "0x0000FFFF",
			"Event":          event,
			"Type":           typ,
		},
	}, &result)
	if err != nil {
		return nil, err
	}

	return result.OPFileQuery, err
}
