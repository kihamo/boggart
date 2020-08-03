package miio

import (
	"context"
	"errors"
	"io"
	"net"
	"time"

	"github.com/kihamo/boggart/providers/xiaomi/miio/internal"
	"github.com/kihamo/boggart/types"
)

const (
	OTAStatusUnknown     otaStatus = "unknown"
	OTAStatusDownloading otaStatus = "downloading"
	OTAStatusInstalling  otaStatus = "installing"
	OTAStatusFailed      otaStatus = "failed"
	OTAStatusIdle        otaStatus = "idle"
)

type otaStatus string

type Device struct {
	io.Closer

	client *Client
}

type InfoPayload struct {
	HardwareVersion string             `json:"hw_ver"`
	FirmwareVersion string             `json:"fw_ver"`
	Token           string             `json:"token"`
	LifeTime        time.Duration      `json:"life"`
	MAC             types.HardwareAddr `json:"mac"`
	Model           string             `json:"model"`
	AccessPoint     struct {
		BSSID types.HardwareAddr `json:"bssid"`
		RSSI  int                `json:"rssi"`
		SSID  string             `json:"ssid"`
	} `json:"ap"`
	Network struct {
		Gateway net.IP       `json:"gw"`
		LocalIP net.IP       `json:"localIp"`
		Mask    types.IPMask `json:"mask"`
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

func (d *Device) ID(ctx context.Context) (uint16, error) {
	return d.client.DeviceID(ctx)
}

func (d *Device) Close() error {
	return d.client.Close()
}

func (d *Device) Info(ctx context.Context) (InfoPayload, error) {
	var response struct {
		Response
		Result InfoPayload `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "miIO.info", nil, &response)
	if err != nil {
		return InfoPayload{}, err
	}

	response.Result.LifeTime *= time.Second

	return response.Result, nil
}

func (d *Device) WiFiStatus(ctx context.Context) (WiFiStatusPayload, error) {
	var response struct {
		Response

		Result WiFiStatusPayload `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "miIO.wifi_assoc_state", nil, &response)
	if err != nil {
		return WiFiStatusPayload{}, err
	}

	return response.Result, nil
}

// проверить загрузку на устройстве можно в /mnt/data/.temp
func (d *Device) OTALocalServer(ctx context.Context, file io.ReadSeeker) error {
	server, err := internal.NewServer(file, d.HostnameForLocalServer())
	if err != nil {
		return err
	}
	defer server.Close()

	err = d.OTA(ctx, server.URL().String(), server.MD5())
	if err == nil {
		ticker := time.NewTicker(time.Second)

		for range ticker.C {
			if status, err := d.otaStatus(ctx); err == nil {
				switch status {
				case OTAStatusFailed:
					return errors.New("OTA install failed")

				case OTAStatusInstalling:
					return nil
				}
			}
		}
	}

	return err
}

func (d *Device) OTA(ctx context.Context, url, md5sum string) error {
	var response ResponseOK

	err := d.Client().CallRPC(ctx, "miIO.ota", map[string]interface{}{
		"mode":     "normal",
		"install":  "1",
		"app_url":  url,
		"file_md5": md5sum,
		"proc":     "dnld install",
	}, &response)
	if err != nil {
		return err
	}

	if !ResponseIsOK(response) {
		return errors.New("device return not OK response")
	}

	return nil
}

func (d *Device) otaStatus(ctx context.Context) (otaStatus, error) {
	var response struct {
		Response
		Result []string `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "miIO.get_ota_state", nil, &response)
	if err != nil {
		return OTAStatusUnknown, err
	}

	return otaStatus(response.Result[0]), nil
}

func (d *Device) OTAProgress(ctx context.Context) (uint64, error) {
	var response struct {
		Response
		Result []uint64 `json:"result"`
	}

	err := d.Client().CallRPC(ctx, "miIO.get_ota_progress", nil, &response)
	if err != nil {
		return 0, err
	}

	return response.Result[0], nil
}

func (d *Device) HostnameForLocalServer() string {
	if addr, err := d.client.LocalAddr(); err == nil {
		return addr.IP.String()
	}

	return ""
}
