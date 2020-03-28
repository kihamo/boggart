package esp

import (
	"sync"

	"github.com/kihamo/boggart/atomic"
)

const (
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
