package hikvision

import (
	"bufio"
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	connectionAttemptDuration = time.Minute

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
	connection Connection
	request    *http.Request
	alerts     chan *EventNotificationAlertStreamResponse
	errors     chan error
}

func NewAlertStreaming(client Connection, request *http.Request) *AlertStreaming {
	return &AlertStreaming{
		connection: client,
		request:    request,
	}
}

func (s *AlertStreaming) Start() {
	s.alerts = make(chan *EventNotificationAlertStreamResponse)
	s.errors = make(chan error)

	go func() {
		ticker := time.NewTicker(connectionAttemptDuration)
		for ; true; <-ticker.C {
			response, err := s.connection.Do(s.request)
			if err != nil {
				s.errors <- err
				continue
			}

			if response.StatusCode != http.StatusOK {
				s.errors <- fmt.Errorf("Response status %d isn't OK", response.StatusCode, response)
				continue
			}

			go func() {
				s.read(response)
			}()

			break
		}
	}()
}

func (s *AlertStreaming) read(response *http.Response) {
	reader := bufio.NewReader(response.Body)
	buf := bytes.NewBuffer(nil)

	for {
		select {
		case <-s.request.Context().Done():
			buf.Reset()
			close(s.alerts)
			close(s.errors)
			return

		default:
			line, err := reader.ReadBytes('\n')
			if err == io.EOF || len(line) == 0 {
				continue
			}

			if err != nil {
				s.errors <- err
				continue
			}

			if bytes.HasPrefix(line, []byte("<EventNotificationAlert")) {
				buf.Write(line)
				continue
			} else if buf.Len() == 0 { // если сообщение не началось игнорируем весь контент
				continue
			}

			buf.Write(line)

			// если сообщение заканчивается запускаем алгоритм
			if !bytes.HasPrefix(line, []byte("</EventNotificationAlert>")) {
				continue
			}

			event := &EventNotificationAlertStreamResponse{}
			err = xml.Unmarshal(buf.Bytes(), event)
			buf.Reset()

			if err != nil {
				s.errors <- err
			} else {
				s.alerts <- event
			}
		}
	}
}

func (s *AlertStreaming) NextAlert() <-chan *EventNotificationAlertStreamResponse {
	return s.alerts
}

func (s *AlertStreaming) NextError() <-chan error {
	return s.errors
}

func (a *ISAPI) EventNotificationAlertStream(ctx context.Context) (*AlertStreaming, error) {
	request, err := http.NewRequest(http.MethodGet, a.address+"/Event/notification/alertStream", nil)
	if err != nil {
		return nil, err
	}

	stream := NewAlertStreaming(a, request.WithContext(ctx))
	stream.Start()

	return stream, nil
}
