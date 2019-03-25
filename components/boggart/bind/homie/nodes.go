package homie

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/kihamo/boggart/components/boggart/atomic"

	"github.com/kihamo/boggart/components/mqtt"
)

const (
	nodesTopicNodes    = MQTTPrefix + "+/+"
	nodesTopicProperty = MQTTPrefix + "+/+/+"

	dataTypeInteger = "integer"
	dataTypeFloat   = "float"
	dataTypeBoolean = "boolean"
	dataTypeString  = "string"
	dataTypeEnum    = "enum"
	dataTypeColor   = "color"
)

type node struct {
	ID         *atomic.String
	Name       *atomic.String
	Type       *atomic.String
	Array      *atomic.String
	properties *sync.Map
}

type nodeProperty struct {
	Name     *atomic.String
	Settable *atomic.Bool
	Retained *atomic.Bool
	Unit     *atomic.String
	DataType *atomic.String
	Format   *atomic.String
	Value    *atomic.String
}

func (n *node) propertyLoadOrStore(name string) *nodeProperty {
	if value, ok := n.properties.Load(name); ok {
		return value.(*nodeProperty)
	}

	property := &nodeProperty{
		Name:     atomic.NewStringDefault(name),
		Settable: atomic.NewBool(),
		Retained: atomic.NewBool(),
		Unit:     atomic.NewString(),
		DataType: atomic.NewString(),
		Format:   atomic.NewString(),
		Value:    atomic.NewString(),
	}

	n.properties.Store(name, property)
	return property
}

func (n *node) Properties() []*nodeProperty {
	result := make([]*nodeProperty, 0)

	n.properties.Range(func(key, value interface{}) bool {
		result = append(result, value.(*nodeProperty))
		return true
	})

	return result
}

func (b *Bind) Nodes() []*node {
	result := make([]*node, 0)

	b.nodes.Range(func(key, value interface{}) bool {
		result = append(result, value.(*node))
		return true
	})

	return result
}

func (b *Bind) nodesLoadOrStore(name string) *node {
	if value, ok := b.nodes.Load(name); ok {
		return value.(*node)
	}

	n := &node{
		ID:         atomic.NewStringDefault(name),
		Name:       atomic.NewString(),
		Type:       atomic.NewString(),
		Array:      atomic.NewString(),
		properties: &sync.Map{},
	}

	b.nodes.Store(name, n)
	return n
}

/*
	homie / device ID / node ID / $node-attribute
	homie / device ID / node ID / property ID
*/
func (b *Bind) nodesSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := mqtt.RouteSplit(message.Topic())

	// skip $fw $stats $implementation etc.
	if strings.HasPrefix(route[2], "$") {
		return nil
	}

	n := b.nodesLoadOrStore(route[2])

	if strings.HasPrefix(route[3], "$") { // node attribute
		switch route[3] {
		case "$name":
			n.Name.Set(message.String())
		case "$type":
			n.Type.Set(message.String())
		case "$properties":
			// TODO
		case "$array":
			// TODO
		}
	} else { // node property value
		property := n.propertyLoadOrStore(route[3])
		property.Value.Set(message.String())
	}

	return nil
}

/*
	homie / device ID / node ID / property ID / $property-attribute
*/
func (b *Bind) nodesPropertySubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := mqtt.RouteSplit(message.Topic())

	// skip $fw $stats $implementation etc.
	if strings.HasPrefix(route[2], "$") {
		return nil
	}

	n := b.nodesLoadOrStore(route[2])
	property := n.propertyLoadOrStore(route[3])

	switch route[4] {
	case "name", "$name":
		property.Name.Set(message.String())
	case "settable", "$settable":
		property.Settable.Set(message.Bool())
	case "retained", "$retained":
		property.Retained.Set(message.Bool())
	case "unit", "$unit":
		property.Unit.Set(message.String())
	case "datatype", "$datatype":
		switch message.String() {
		case dataTypeInteger, dataTypeFloat, dataTypeBoolean, dataTypeString, dataTypeEnum, dataTypeColor:
			property.DataType.Set(message.String())
		default:
			return errors.New("unknown property type " + message.String())
		}
	case "format", "$format":
		property.Format.Set(message.String())
	}

	return nil
}
