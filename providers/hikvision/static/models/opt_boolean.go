package models

import (
	"encoding/xml"
	"strconv"
	"strings"
)

type OptBoolean struct {
	Opt

	opts  []bool
	value bool
}

func (m *OptBoolean) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	if err := d.DecodeElement(&m.value, &start); err != nil {
		return err
	}

	for _, a := range start.Attr {
		if a.Name.Local == "opt" {
			m.opts = make([]bool, 0, 2)

			for _, f := range strings.Split(a.Value, ",") {
				if opt, err := strconv.ParseBool(f); err != nil {
					return err
				} else {
					m.opts = append(m.opts, opt)
				}
			}

			break
		}
	}

	return nil
}

func (m *OptBoolean) Value() bool {
	return m.value
}

func (m *OptBoolean) Options() []bool {
	return m.opts
}
