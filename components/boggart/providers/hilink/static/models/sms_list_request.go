package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type SMSListRequest struct {
	// Важен порядок полей, иначе метод возвращает ошибку
	PageIndex       int64 `json:"PageIndex,omitempty" xml:"PageIndex"`
	ReadCount       int64 `json:"ReadCount,omitempty" xml:"ReadCount"`
	BoxType         int64 `json:"BoxType,omitempty" xml:"BoxType"`
	SortType        int64 `json:"SortType,omitempty" xml:"SortType"`
	Ascending       int64 `json:"Ascending,omitempty" xml:"Ascending"`
	UnreadPreferred int64 `json:"UnreadPreferred,omitempty" xml:"UnreadPreferred"`
}

func (m *SMSListRequest) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *SMSListRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

func (m *SMSListRequest) UnmarshalBinary(b []byte) error {
	var res SMSListRequest
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
