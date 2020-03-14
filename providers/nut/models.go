package nut

import (
	"strconv"
	"strings"
)

type ListUPS struct {
	Name        string
	Description string
}

type ListVariable struct {
	UPS   string
	Name  string
	Value string
}

type ListCommand struct {
	UPS  string
	Name string
}

type Type struct {
	Writeable bool
	MaxLength int
	Name      string
}

func (t Type) ConvertValue(value string) interface{} {
	if t.Name == "NUMBER" {
		if strings.Contains(value, ".") {
			if v, err := strconv.ParseFloat(value, 64); err == nil {
				return v
			}
		} else {
			if v, err := strconv.Atoi(value); err == nil {
				return v
			}
		}
	}

	switch value {
	case "enabled":
		return true

	case "disabled":
		return false
	}

	return value
}
