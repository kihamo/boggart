package hikvision

import (
	"context"
	"net/http"

	"github.com/kihamo/shadow/components/tracing"
)

const (
	contentManagerPrefixURL = "/ContentMgmt"
)

type ContentManagementStorageResponse struct {
	HDD []struct {
		ID        uint64 `xml:"id"`
		Name      string `xml:"hddName"`
		Path      string `xml:"hddPath"`
		Type      string `xml:"hddType"`
		Status    string `xml:"status"`
		Capacity  uint64 `xml:"capacity"`
		FreeSpace uint64 `xml:"freeSpace"`
		Property  string `xml:"property"`
	} `xml:"hddList>hdd"`
	NAS []struct {
		ID                   uint64 `xml:"id"`
		AddressingFormatType string `xml:"addressingFormatType"`
	} `xml:"nasList>nas"`
}

func (a *ISAPI) ContentManagementStorage(ctx context.Context) (result ContentManagementStorageResponse, err error) {
	span, ctx := tracing.StartSpanFromContext(ctx, ComponentName, "content_management.storage")
	defer span.Finish()

	err = a.DoXML(ctx, http.MethodGet, a.address+contentManagerPrefixURL+"/Storage", nil, &result)
	if err != nil {
		tracing.SpanError(span, err)
	}

	return result, err
}
