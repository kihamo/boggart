package hikvision

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/kihamo/shadow/components/tracing"
)

const (
	systemPrefixURL = "/System"
)

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
