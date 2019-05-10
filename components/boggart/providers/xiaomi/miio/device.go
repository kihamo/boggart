package miio

import (
	"context"
	"io"
	"net"
	"time"

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
	LifeTime        time.Duration        `json:"life"`
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

type WiFiStatusPayload struct {
	State               string `json:"state"`
	AuthFailCount       uint64 `json:"auth_fail_count"`
	ConnectSuccessCount uint64 `json:"conn_success_count"`
	ConnectFailCount    uint64 `json:"conn_fail_count"`
	DHCPFailCount       uint64 `json:"dhcp_fail_count"`
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

func (d *Device) Info(ctx context.Context) (InfoPayload, error) {
	var reply struct {
		Response
		Result InfoPayload `json:"result"`
	}

	err := d.Client().Send(ctx, "miIO.info", nil, &reply)
	if err != nil {
		return InfoPayload{}, err
	}

	reply.Result.LifeTime *= time.Second

	return reply.Result, nil
}

func (d *Device) WiFiStatus(ctx context.Context) (WiFiStatusPayload, error) {
	type response struct {
		Response

		Result WiFiStatusPayload `json:"result"`
	}

	var reply response

	err := d.Client().Send(ctx, "miIO.wifi_assoc_state", nil, &reply)
	if err != nil {
		return WiFiStatusPayload{}, err
	}

	return reply.Result, nil
}

/*
func (d *Device) OTAStatus(ctx context.Context) (error) {
	var reply struct {
		Response
		// Result []uint64 `json:"result"`
	}

	err := d.Client().Send(ctx, "miIO.get_ota_state", nil, &reply)
	if err != nil {
		return err
	}

	fmt.Println(reply.Result)

	return nil
}
*/

func (d *Device) OTAProgress(ctx context.Context) (uint64, error) {
	var reply struct {
		Response
		Result []uint64 `json:"result"`
	}

	err := d.Client().Send(ctx, "miIO.get_ota_progress", nil, &reply)
	if err != nil {
		return 0, err
	}

	return reply.Result[0], nil
}
