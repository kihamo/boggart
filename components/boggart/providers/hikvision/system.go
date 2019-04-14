package hikvision

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kihamo/shadow/components/tracing"
)

const (
	systemPrefixURL = "/System"
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

type SystemUpgradeStatusResponse struct {
	Upgrading bool   `xml:"upgrading"`
	Percent   uint64 `xml:"percent"`
}

func (a *ISAPI) SystemDeviceInfo(ctx context.Context) (result SystemDeviceInfoResponse, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "system.device_info")
	defer span.Finish()

	err = a.DoXML(ctx, http.MethodGet, a.address+systemPrefixURL+"/deviceInfo", nil, &result)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return result, err
}

func (a *ISAPI) SystemStatus(ctx context.Context) (result SystemStatusResponse, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "system.status")
	defer span.Finish()

	err = a.DoXML(ctx, http.MethodGet, a.address+systemPrefixURL+"/status", nil, &result)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return result, err
}

func (a *ISAPI) SystemUpdateFirmware(ctx context.Context, firmware io.Reader) (status ResponseStatus, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "system.update_firmware")

	defer func() {
		if err != nil {
			tracing.SpanError(span, err)
		}

		span.Finish()
	}()

	// протокол требует обязательного вычисления длины входящего пакета
	buf := bytes.NewBuffer(nil)
	contentLength, err := io.Copy(buf, firmware)
	if err != nil {
		return status, err
	}

	request, err := http.NewRequest(http.MethodPut, a.address+systemPrefixURL+"/updateFirmware", buf)
	if err != nil {
		return status, err
	}

	request.ContentLength = contentLength
	request.Header.Set("Content-Type", `application/x-www-form-urlencoded`)

	response, err := a.DoRequest(ctx, request)

	defer response.Body.Close()

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return status, err
	}

	err = xml.Unmarshal(content, &status)

	return status, err
}

func (a *ISAPI) SystemUpgradeStatus(ctx context.Context) (result SystemUpgradeStatusResponse, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "system.upgrade_status")
	defer span.Finish()

	err = a.DoXML(ctx, http.MethodGet, a.address+systemPrefixURL+"/upgradeStatus", nil, &result)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return result, err
}

func (a *ISAPI) SystemReboot(ctx context.Context) (err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "system.reboot")
	defer span.Finish()

	_, err = a.Do(ctx, http.MethodPut, a.address+systemPrefixURL+"/reboot", nil)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return err
}
