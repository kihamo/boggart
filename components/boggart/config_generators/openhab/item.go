package openhab

import (
	"strings"
)

const (
	ItemTypeColor         = "Color"
	ItemTypeContact       = "Contact"
	ItemTypeDateTime      = "DateTime"
	ItemTypeDimmer        = "Dimmer"
	ItemTypeGroup         = "Group"
	ItemTypeImage         = "Image"
	ItemTypeLocation      = "Location"
	ItemTypeNumber        = "Number"
	ItemTypePlayer        = "Player"
	ItemTypeRollerShutter = "RollerShutter"
	ItemTypeString        = "String"
	ItemTypeSwitch        = "Switch"
)

type Items []*Item

func (i Items) String() string {
	var s string

	for i, item := range i {
		if i > 0 {
			s += "\n"
		}

		s += item.String()
	}

	return s
}

type Item struct {
	name       string
	typ        string
	label      string
	icon       string
	groups     []string
	tags       []string
	parameters *Parameters
}

func NewItem(name, typ string) *Item {
	return &Item{
		name:       name,
		typ:        typ,
		parameters: newParameters(),
		groups:     make([]string, 0),
		tags:       make([]string, 0),
	}
}

func (i *Item) WithLabel(label string) *Item {
	i.label = label
	return i
}

func (i *Item) WithIcon(icon string) *Item {
	i.icon = icon
	return i
}

func (i *Item) WithGroups(groups ...string) *Item {
	i.groups = groups
	return i
}

func (i *Item) WithTags(tags ...string) *Item {
	i.tags = tags
	return i
}

func (i *Item) WithChannel(channel *Channel) *Item {
	const key = "channel"

	if channel == nil {
		i.parameters.Delete(key)
		return i
	}

	i.parameters.Set(key, channel.ChannelID())
	return i
}

func (i *Item) WithParameter(key string, value interface{}) *Item {
	i.parameters.Set(key, value)
	return i
}

func (i *Item) String() string {
	s := i.typ + " " + i.name

	if i.label != "" {
		s += " \"" + i.label + "\""
	}

	if i.icon != "" {
		s += " <" + i.icon + ">"
	}

	if len(i.groups) > 0 {
		s += " (" + strings.Join(i.groups, ",") + ")"
	}

	if len(i.tags) > 0 {
		s += " [\"" + strings.Join(i.tags, "\",\"") + "\"]"
	}

	if i.parameters != nil {
		if params := i.parameters.String(); params != "" {
			s += " {" + params + "}"
		}
	}

	return s
}
