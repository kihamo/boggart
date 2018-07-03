package hikvision

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

const (
	boundary = "boundary"

	EventTypeIO                   = "IO"
	EventTypeVMD                  = "VMD"
	EventTypeVideoLoss            = "videoloss"
	EventTypeShelterAlarm         = "shelteralarm"
	EventTypeFaceDetection        = "facedetection"
	EventTypeDefocus              = "defocus"
	EventTypeAudioException       = "audioexception"
	EventTypeSceneChangeDetection = "scenechangedetection"
	EventTypeFieldDetection       = "fielddetection"
	EventTypeLineDetection        = "linedetection"
	EventTypeRegionEntrance       = "regionEntrance"
	EventTypeRegionExiting        = "regionExiting"
	EventTypeLoitering            = "loitering"
	EventTypeGroup                = "group"
	EventTypeRapidMove            = "rapidMove"
	EventTypeParking              = "parking"
	EventTypeUnattendedBaggage    = "unattendedBaggage"
	EventTypeAttendedBaggage      = "attendedBaggage"
	EventTypePIR                  = "PIR"
	EventTypePeopleDetection      = "peopleDetection"

	EventEventStateActive  = "active"
	EventEventStateInctive = "inactive"
)

type EventNotificationAlertStreamResponse struct {
	IpAddress           string    `xml:"ipAddress"`
	Ipv6Address         string    `xml:"ipv6Address"`
	PortNo              uint64    `xml:"portNo"`
	Protocol            string    `xml:"protocol"`
	MacAddress          string    `xml:"macAddress"`
	DynChannelID        uint64    `xml:"dynChannelID"`
	DateTime            time.Time `xml:"dateTime"`
	ActivePostCount     uint64    `xml:"activePostCount"`
	EventType           string    `xml:"eventType"`
	EventState          string    `xml:"eventState"`
	EventDescription    string    `xml:"eventDescription"`
	InputIOPortID       uint64    `xml:"inputIOPortID"`
	DynInputIOPortID    string    `xml:"dynInputIOPortID"`
	DetectionRegionList []struct {
		RegionID         string `xml:"regionID"`
		SensitivityLevel uint64 `xml:"sensitivityLevel"`
	} `xml:"DetectionRegionList>DetectionRegionEntry"`
}

type AlertStreaming struct {
	reader *bufio.Reader
}

func (s *AlertStreaming) NextAlert() (*EventNotificationAlertStreamResponse, error) {
	buf := bytes.NewBuffer([]byte("--" + boundary + "\n"))

	var (
		isStart bool
		event   *EventNotificationAlertStreamResponse
	)

	for {
		line, err := s.reader.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		if bytes.HasPrefix(line, []byte("--"+boundary)) {
			if isStart {
				reader := multipart.NewReader(buf, boundary)

				part, err := reader.NextPart()
				if err != nil {
					return nil, err
				}

				content, _ := ioutil.ReadAll(part)
				event = &EventNotificationAlertStreamResponse{}

				if err := xml.Unmarshal(content, event); err != nil {
					return nil, err
				}

				return event, nil
			}
		} else {
			buf.Write(line)
			isStart = true
		}
	}
}

func (a *ISAPI) EventNotificationAlertStream(ctx context.Context) (*AlertStreaming, error) {
	request, err := http.NewRequest(http.MethodGet, a.address+"/Event/notification/alertStream", nil)
	if err != nil {
		return nil, err
	}

	response, err := a.do(request.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	return &AlertStreaming{
		reader: bufio.NewReader(response.Body),
	}, nil
}
