package openhab

import (
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	ChannelFormatBeforePublish = "%s"
	ChannelOnBrightness        = 10

	ChannelTypeString        = "string"
	ChannelTypeNumber        = "number"
	ChannelTypeDimmer        = "dimmer"
	ChannelTypeContact       = "contact"
	ChannelTypeSwitch        = "switch"
	ChannelTypeColor         = "color"
	ChannelTypeLocation      = "location"
	ChannelTypeImage         = "image"
	ChannelTypeDateTime      = "datetime"
	ChannelTypeRollerShutter = "rollershutter"

	ColorModeHSB = "hsb"
	ColorModeRGB = "rgb"
	ColorModeXYY = "xyY"
)

type Channel struct {
	typ            string
	id             string
	genericThingID string
	label          string
	parameters     *Parameters
	items          Items
}

func NewChannel(id, typ string) *Channel {
	// TODO: id не может называться State (экспериментально полученное значение)

	return (&Channel{
		id:         IDNormalize(id),
		typ:        typ,
		parameters: newParameters(),
		items:      make(Items, 0),
	}).
		WithFormatBeforePublish(ChannelFormatBeforePublish).
		WithOnBrightness(ChannelOnBrightness)
}

func (c *Channel) ChannelID() string {
	return c.genericThingID + ":" + c.id
}

func (c *Channel) WithGenericThing(thing *GenericThing) *Channel {
	c.genericThingID = thing.GenericThingID()

	for _, item := range c.items {
		item.WithChannel(c)
	}

	return c
}

func (c *Channel) WithLabel(label string) *Channel {
	c.label = label
	return c
}

func (c *Channel) WithStateTopic(stateTopic mqtt.Topic) *Channel {
	const key = "stateTopic"

	if stateTopic == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, stateTopic)
	return c
}

func (c *Channel) WithTransformationPattern(transformationPattern string) *Channel {
	const key = "transformationPattern"

	if transformationPattern == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, transformationPattern)
	return c
}

func (c *Channel) WithTransformationPatternOut(transformationPatternOut string) *Channel {
	const key = "transformationPatternOut"

	if transformationPatternOut == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, transformationPatternOut)
	return c
}

func (c *Channel) WithCommandTopic(commandTopic mqtt.Topic) *Channel {
	const key = "commandTopic"

	if commandTopic == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, commandTopic)
	return c
}

func (c *Channel) WithFormatBeforePublish(formatBeforePublish string) *Channel {
	const key = "formatBeforePublish"

	if formatBeforePublish == ChannelFormatBeforePublish {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, formatBeforePublish)
	return c
}

func (c *Channel) WithPostCommand(postCommand bool) *Channel {
	const key = "postCommand"

	if !postCommand {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, postCommand)
	return c
}

func (c *Channel) WithRetained(retained bool) *Channel {
	const key = "retained"

	if !retained {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, retained)
	return c
}

func (c *Channel) WithQOS(qos int) *Channel {
	const key = "qos"

	if qos <= 0 || qos > 2 {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, qos)
	return c
}

func (c *Channel) WithTrigger(trigger bool) *Channel {
	const key = "trigger"

	if !trigger {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, trigger)
	return c
}

func (c *Channel) WithAllowedStates(allowedStates []string) *Channel {
	if c.typ != ChannelTypeString {
		return c
	}

	const key = "allowedStates"

	if len(allowedStates) == 0 {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, allowedStates)
	return c
}

func (c *Channel) WithMin(min float64) *Channel {
	if c.typ != ChannelTypeNumber && c.typ != ChannelTypeDimmer {
		return c
	}

	const key = "min"

	if min == 0 {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, min)
	return c
}

func (c *Channel) WithMax(max float64) *Channel {
	if c.typ != ChannelTypeNumber && c.typ != ChannelTypeDimmer {
		return c
	}

	const key = "max"

	if max == 0 {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, max)
	return c
}

func (c *Channel) WithStep(step float64) *Channel {
	if c.typ != ChannelTypeNumber && c.typ != ChannelTypeDimmer {
		return c
	}

	const key = "step"

	if step == 0 {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, step)
	return c
}

func (c *Channel) WithUnit(unit string) *Channel {
	if c.typ != ChannelTypeNumber {
		return c
	}

	const key = "unit"

	if unit == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, unit)
	return c
}

func (c *Channel) WithOn(on string) *Channel {
	if c.typ != ChannelTypeDimmer && c.typ != ChannelTypeContact && c.typ != ChannelTypeSwitch && c.typ != ChannelTypeColor && c.typ != ChannelTypeRollerShutter {
		return c
	}

	const key = "on"

	if on == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, on)
	return c
}

func (c *Channel) WithOff(off string) *Channel {
	if c.typ != ChannelTypeDimmer && c.typ != ChannelTypeContact && c.typ != ChannelTypeSwitch && c.typ != ChannelTypeColor && c.typ != ChannelTypeRollerShutter {
		return c
	}

	const key = "off"

	if off == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, off)
	return c
}

func (c *Channel) WithColorMode(colorMode string) *Channel {
	if c.typ != ChannelTypeColor {
		return c
	}

	const key = "color_mode"

	if colorMode == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, colorMode)
	return c
}

func (c *Channel) WithOnBrightness(onBrightness int) *Channel {
	if c.typ != ChannelTypeColor {
		return c
	}

	const key = "onBrightness"

	if onBrightness == ChannelOnBrightness {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, onBrightness)
	return c
}

func (c *Channel) WithStop(stop string) *Channel {
	if c.typ != ChannelTypeRollerShutter {
		return c
	}

	const key = "stop"

	if stop == "" {
		c.parameters.Delete(key)
		return c
	}

	c.parameters.Set(key, stop)
	return c
}

func (c *Channel) AddItems(items ...*Item) *Channel {
	for _, item := range items {
		item.WithChannel(c)
		c.items = append(c.items, item)
	}

	return c
}

func (c *Channel) Items() Items {
	return c.items
}

// Type string : Heating "Force" [ stateTopic="zigbee2mqtt/radiator_bedroom", transformationPattern="JSONPATH:$.force"]
func (c *Channel) String() string {
	s := "Type " + c.typ + " : " + c.id

	if c.label != "" {
		s += " \"" + c.label + "\""
	}

	if c.parameters != nil {
		if params := c.parameters.String(); params != "" {
			s += " [" + params + "]"
		}
	}

	return s
}
