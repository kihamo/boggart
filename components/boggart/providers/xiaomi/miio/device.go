package miio

import (
	"io"
	"net"

	"github.com/kihamo/boggart/components/boggart"
)

type Device struct {
	io.Closer

	client *Client
}

type InfoPayload struct {
	HardwareVersion string               `json:"hw_ver"`
	FirmwareVersion string               `json:"fw_ver"`
	Token           string               `json:"token"`
	LifeTime        uint64               `json:"life"`
	MAC             boggart.HardwareAddr `json:"mac"`
	Model           string               `json:"model"`
	AccessPoint     struct {
		BSSID boggart.HardwareAddr `json:"bssid"`
		RSSI  int                  `json:"rssi"`
		SSID  string               `json:"ssid"`
	} `json:"ap"`
	Network struct {
		Gateway net.IP         `json:"gw"`
		LocalIP net.IP         `json:"localIp"`
		Mask    boggart.IPMask `json:"mask"`
	} `json:"netif"`
}

func NewDevice(address, token string) *Device {
	return &Device{
		client: NewClient(address, token),
	}
}

func (d *Device) Client() *Client {
	return d.client
}

func (d *Device) ID() ([]byte, error) {
	return d.client.DeviceID()
}

func (d *Device) Close() error {
	return d.client.Close()
}

func (d *Device) Info() (InfoPayload, error) {
	var reply struct {
		Response
		Result InfoPayload `json:"result"`
	}

	err := d.Client().Send("miIO.info", nil, &reply)
	if err != nil {
		return InfoPayload{}, err
	}

	return reply.Result, nil
}

/*
func (d *Device) OTAStatus() (error) {
	var reply struct {
		Response
		// Result []uint64 `json:"result"`
	}

	err := d.Client().Send("miIO.get_ota_state", nil, &reply)
	if err != nil {
		return err
	}

	fmt.Println(reply.Result)

	return nil
}
*/

func (d *Device) OTAProgress() (uint64, error) {
	var reply struct {
		Response
		Result []uint64 `json:"result"`
	}

	err := d.Client().Send("miIO.get_ota_progress", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}
