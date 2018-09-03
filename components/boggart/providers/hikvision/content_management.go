package hikvision

import (
	"context"
	"net/http"
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

func (a *ISAPI) ContentManagementStorage(ctx context.Context) (ContentManagementStorageResponse, error) {
	result := ContentManagementStorageResponse{}

	request, err := http.NewRequest(http.MethodGet, a.address+"/ContentMgmt/Storage", nil)
	if err != nil {
		return result, err
	}

	request = request.WithContext(ctx)
	err = a.DoAndParse(request, &result)

	return result, err
}
