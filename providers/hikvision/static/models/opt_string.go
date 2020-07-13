package models

import (
	"encoding/xml"
	"strings"
)

type OptString struct {
	Opt

	opts  []string
	value string
}

func (m *OptString) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if err := d.DecodeElement(&m.value, &start); err != nil {
		return err
	}

	for _, a := range start.Attr {
		if a.Name.Local == "opt" {
			m.opts = make([]string, 0)

			for _, opt := range strings.Split(a.Value, ",") {
				m.opts = append(m.opts, opt)
			}

			break
		}
	}

	return nil
}

func (m *OptString) Value() string {
	return m.value
}

func (m *OptString) Options() []string {
	return m.opts
}
