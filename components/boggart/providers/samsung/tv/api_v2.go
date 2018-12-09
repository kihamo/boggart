package tv

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kihamo/boggart/components/boggart/protocols/http"
)

const (
	ApiV2BasePath      = "/api/v2/"
	ApiV2LegacyPort    = 55000
	ApiV2WebSocketPort = 8001
	ApiV2ApplicationId = "SamsungTV"
)

const (
	KeyPower = "KEY_POWER"
)

type ApiV2DeviceResponse struct {
	Device struct {
		OS                string
		CountryCode       string
		Description       string
		DeveloperIP       string
		DeveloperMode     string
		UID               string
		FirmwareVersion   string
		ID                string
		IP                string
		Model             string
		ModelName         string
		Name              string
		NetworkType       string
		Resolution        string
		Type              string
		UDN               string
		WifiMac           string
		FrameTVSupport    bool
		GamePadSupport    bool
		ImeSyncedSupport  bool
		VoiceSupport      bool
		SmartHubAgreement bool
	}
	ID      string
	Name    string
	Remote  string
	Type    string
	Version string
	Support struct {
		DRMPlayReady         bool
		DRMWideVine          bool
		DMPAvailable         bool
		EDENAvailable        bool
		FrameTVSupport       bool
		ImeSyncedSupport     bool
		RemoteAvailable      bool
		RemoteFourDirections bool
		RemoteTouchPad       bool
		RemoteVoiceControl   bool
	}
}

func (r *ApiV2DeviceResponse) UnmarshalJSON(b []byte) error {
	var t map[string]interface{}

	if err := json.Unmarshal(b, &t); err != nil {
		return err
	}

	r.ID = t["id"].(string)
	r.Name = t["name"].(string)
	r.Remote = t["remote"].(string)
	r.Type = t["type"].(string)
	r.Version = t["version"].(string)

	if device, ok := t["device"]; ok {
		if deviceMap, ok := device.(map[string]interface{}); ok {
			r.Device.OS = deviceMap["OS"].(string)
			r.Device.CountryCode = deviceMap["countryCode"].(string)
			r.Device.Description = deviceMap["description"].(string)
			r.Device.DeveloperIP = deviceMap["developerIP"].(string)
			r.Device.DeveloperMode = deviceMap["developerMode"].(string)
			r.Device.UID = deviceMap["duid"].(string)
			r.Device.FirmwareVersion = deviceMap["firmwareVersion"].(string)
			r.Device.ID = deviceMap["id"].(string)
			r.Device.IP = deviceMap["ip"].(string)
			r.Device.Model = deviceMap["model"].(string)
			r.Device.ModelName = deviceMap["modelName"].(string)
			r.Device.Name = deviceMap["name"].(string)
			r.Device.NetworkType = deviceMap["networkType"].(string)
			r.Device.Resolution = deviceMap["resolution"].(string)
			r.Device.Type = deviceMap["type"].(string)
			r.Device.UDN = deviceMap["udn"].(string)
			r.Device.WifiMac = deviceMap["wifiMac"].(string)

			if flag, ok := deviceMap["FrameTVSupport"]; ok {
				r.Device.FrameTVSupport = flag == "true"
			}

			if flag, ok := deviceMap["GamePadSupport"]; ok {
				r.Device.GamePadSupport = flag == "true"
			}

			if flag, ok := deviceMap["ImeSyncedSupport"]; ok {
				r.Device.ImeSyncedSupport = flag == "true"
			}

			if flag, ok := deviceMap["VoiceSupport"]; ok {
				r.Device.VoiceSupport = flag == "true"
			}

			if flag, ok := deviceMap["smartHubAgreement"]; ok {
				r.Device.SmartHubAgreement = flag == "true"
			}
		}
	}

	if isSupport, ok := t["isSupport"]; ok {
		var isSupportMap map[string]string

		if err := json.Unmarshal([]byte(isSupport.(string)), &isSupportMap); err != nil {
			return err
		}

		if flag, ok := isSupportMap["DMP_DRM_PLAYREADY"]; ok {
			r.Support.DRMPlayReady = flag == "true"
		}

		if flag, ok := isSupportMap["DMP_DRM_WIDEVINE"]; ok {
			r.Support.DRMWideVine = flag == "true"
		}

		if flag, ok := isSupportMap["DMP_available"]; ok {
			r.Support.DMPAvailable = flag == "true"
		}

		if flag, ok := isSupportMap["EDEN_available"]; ok {
			r.Support.EDENAvailable = flag == "true"
		}

		if flag, ok := isSupportMap["FrameTVSupport"]; ok {
			r.Support.FrameTVSupport = flag == "true"
		}

		if flag, ok := isSupportMap["ImeSyncedSupport"]; ok {
			r.Support.ImeSyncedSupport = flag == "true"
		}

		if flag, ok := isSupportMap["remote_available"]; ok {
			r.Support.RemoteAvailable = flag == "true"
		}

		if flag, ok := isSupportMap["remote_fourDirections"]; ok {
			r.Support.RemoteFourDirections = flag == "true"
		}

		if flag, ok := isSupportMap["remote_touchPad"]; ok {
			r.Support.RemoteTouchPad = flag == "true"
		}

		if flag, ok := isSupportMap["remote_voiceControl"]; ok {
			r.Support.RemoteVoiceControl = flag == "true"
		}
	}

	return nil
}

type ApiV2SendCommandRequest struct {
	Method string `json:"method"`
	Params struct {
		Cmd          string
		DataOfCmd    string
		Option       bool
		TypeOfRemote string
	} `json:"params"`
}

type ApiV2 struct {
	mutex sync.RWMutex

	host    string
	client  *http.Client
	connect *websocket.Conn
	info    *ApiV2DeviceResponse
}

func NewApiV2(host string) *ApiV2 {
	return &ApiV2{
		host:   host,
		client: http.NewClient(),
	}
}

func (a *ApiV2) Host() string {
	return a.host
}

func (a *ApiV2) support() *ApiV2DeviceResponse {
	a.mutex.RLock()
	info := a.info
	a.mutex.RUnlock()

	if info != nil {
		return info
	}

	device, err := a.Device(context.Background())
	if err != nil {
		return &ApiV2DeviceResponse{}
	}

	a.mutex.Lock()
	a.info = &device
	a.mutex.Unlock()

	return &device
}

func (a *ApiV2) RemoteControlConnect() (*websocket.Conn, error) {
	if !a.support().Support.RemoteAvailable {
		return nil, errors.New("remote control isn't supported")
	}

	a.mutex.RLock()
	connect := a.connect
	a.mutex.RUnlock()

	if connect != nil {
		return connect, nil
	}

	dialer := websocket.Dialer{
		HandshakeTimeout: time.Second,
	}

	u := url.URL{
		Scheme:  "ws",
		Host:    net.JoinHostPort(a.host, strconv.Itoa(ApiV2WebSocketPort)),
		Path:    ApiV2BasePath + "channels/samsung.remote.control",
		RawPath: "name=" + base64.StdEncoding.EncodeToString([]byte(ApiV2ApplicationId)),
	}

	connect, _, err := dialer.Dial(u.String(), nil)

	if err == nil {
		_, _, err = connect.ReadMessage()
	}

	if err == nil {
		a.mutex.Lock()
		a.connect = connect
		a.mutex.Unlock()
	}

	return connect, err
}

func (a *ApiV2) Device(ctx context.Context) (ApiV2DeviceResponse, error) {
	reply := ApiV2DeviceResponse{}

	u := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(a.host, strconv.Itoa(ApiV2WebSocketPort)),
		Path:   ApiV2BasePath,
	}

	response, err := a.client.Get(ctx, u.String())
	if err != nil {
		return reply, err
	}

	if err := http.JsonUnmarshal(response, &reply); err != nil {
		return reply, err
	}

	return reply, nil
}

func (a *ApiV2) SendCommand(command string) error {
	connect, err := a.RemoteControlConnect()
	if err != nil {
		return err
	}

	msg := ApiV2SendCommandRequest{
		Method: "ms.remote.control",
		Params: struct {
			Cmd          string
			DataOfCmd    string
			Option       bool
			TypeOfRemote string
		}{
			Cmd:          "Click",
			DataOfCmd:    command,
			Option:       false,
			TypeOfRemote: "SendRemoteKey",
		},
	}

	return connect.WriteJSON(msg)
}
