package xmeye

import (
	"context"
	"time"
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

	MaxLogRecordsPerPage = 128
)

func (c *Client) LogSearch(ctx context.Context, begin, end time.Time, position uint64) ([]LogSearch, error) {
	var result struct {
		Response
		OPLogQuery []LogSearch
	}

	err := c.CallWithResult(ctx, CmdLogSearchRequest, map[string]interface{}{
		"Name":      "OPLogQuery",
		"SessionID": c.sessionIDAsString(),
		"OPLogQuery": map[string]interface{}{
			"BeginTime":   begin.Format(timeLayout),
			"EndTime":     end.Format(timeLayout),
			"LogPosition": position,
			"Type":        "LogAll",
		},
	}, &result)
	if err != nil {
		return nil, err
	}

	return result.OPLogQuery, err
}

func (c *Client) FileSearch(ctx context.Context, begin, end time.Time, channel uint32, event fileSearchEvent, typ fileSearchType) ([]FileSearch, error) {
	var result struct {
		Response
		OPFileQuery []FileSearch
	}

	err := c.CallWithResult(ctx, CmdFileSearchRequest, map[string]interface{}{
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
