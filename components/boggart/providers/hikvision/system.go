package hikvision

import (
	"context"
	"net/http"
	"time"
)

type SystemDeviceInfoResponse struct {
	DeviceName           string `xml:"deviceName"`
	DeviceID             string `xml:"deviceID"`
	DeviceDescription    string `xml:"deviceDescription"`
	DeviceLocation       string `xml:"deviceLocation"`
	SystemContact        string `xml:"systemContact"`
	Model                string `xml:"model"`
	SerialNumber         string `xml:"serialNumber"`
	MacAddress           string `xml:"macAddress"`
	FirmwareVersion      string `xml:"firmwareVersion"`
	FirmwareVersionInfo  string `xml:"firmwareVersionInfo"` // from real device (DS-2DE5220IW-AE)
	FirmwareReleasedDate string `xml:"firmwareReleasedDate"`
	BootVersion          string `xml:"bootVersion"`
	BootReleasedDate     string `xml:"bootReleasedDate"`
	HardwareVersion      string `xml:"hardwareVersion"`
	EncoderVersion       string `xml:"encoderVersion"`
	EncoderReleasedDate  string `xml:"encoderReleasedDate"`
	DecoderVersion       string `xml:"decoderVersion"`
	DecoderReleasedDate  string `xml:"decoderReleasedDate"`
	DeviceType           string `xml:"deviceType"`
	TelecontrolID        uint64 `xml:"telecontrolID"`
	SupportBeep          bool   `xml:"supportBeep"`
	SupportVideoLoss     bool   `xml:"supportVideoLoss"` // from real device (DS-2DE5220IW-AE)
}

type SystemStatusResponse struct {
	CurrentDeviceTime time.Time `xml:"currentDeviceTime"`
	DeviceUpTime      uint64    `xml:"deviceUpTime"`
	Memory            []struct {
		MemoryDescription string          `xml:"memoryDescription"`
		MemoryUsage       overrideFloat64 `xml:"memoryUsage"`
		MemoryAvailable   overrideFloat64 `xml:"memoryAvailable"`
	} `xml:"MemoryList>Memory"`
}

func (a *ISAPI) SystemDeviceInfo(ctx context.Context) (SystemDeviceInfoResponse, error) {
	result := SystemDeviceInfoResponse{}

	request, err := http.NewRequest(http.MethodGet, a.address+"/System/deviceInfo", nil)
	if err != nil {
		return result, err
	}

	request = request.WithContext(ctx)
	err = a.DoAndParse(request, &result)

	return result, err
}

func (a *ISAPI) SystemStatus(ctx context.Context) (SystemStatusResponse, error) {
	result := SystemStatusResponse{}

	request, err := http.NewRequest(http.MethodGet, a.address+"/System/status", nil)
	if err != nil {
		return result, err
	}

	request = request.WithContext(ctx)
	err = a.DoAndParse(request, &result)

	return result, err
}
