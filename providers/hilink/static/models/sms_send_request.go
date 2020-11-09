package models

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

type SMSSendRequest struct {
	// Важен порядок полей, иначе метод возвращает ошибку
	Index    int64    `json:"Index,omitempty" xml:"Index,omitempty"`
	Phones   []string `json:"Phones" xml:"Phones>Phone"`
	Sca      string   `json:"Sca" xml:"Sca"`
	Content  string   `json:"Content,omitempty" xml:"Content,omitempty"`
	Length   int64    `json:"Length,omitempty" xml:"Length,omitempty"`
	Reserved int64    `json:"Reserved,omitempty" xml:"Reserved,omitempty"`
	Date     string   `json:"Date,omitempty" xml:"Date,omitempty"`
}

func (m *SMSSendRequest) Validate(formats strfmt.Registry) error {
	return nil
}

func (m *SMSSendRequest) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}

	return swag.WriteJSON(m)
}

func (m *SMSSendRequest) UnmarshalBinary(b []byte) error {
	var res SMSSendRequest

	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}

	*m = res

	return nil
}
