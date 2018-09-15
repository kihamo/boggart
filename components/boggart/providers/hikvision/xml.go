package hikvision

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// Проблемы с \n в ответе от видео регистратора
type overrideFloat64 struct {
	value float64
}

func (f *overrideFloat64) Value() float64 {
	return f.value
}

func (f *overrideFloat64) Float64() float64 {
	return f.Value()
}

func (f *overrideFloat64) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string
	d.DecodeElement(&raw, &start)

	value, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return err
	}

	*f = overrideFloat64{value}

	return nil
}
