package esp

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/mqtt"
)

const (
	configNameSeparator = "."
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind

	config     *Config
	lastUpdate *atomic.TimeNull

	deviceAttributes sync.Map
	nodes            sync.Map
	settings         sync.Map

	otaEnabled  atomic.Bool
	otaRun      atomic.Bool
	otaWritten  atomic.Uint32
	otaTotal    atomic.Uint32
	otaChecksum atomic.String
	otaFlash    chan struct{}

	status atomic.BoolNull
}

func (b *Bind) updateStatus(status bool) {
	b.status.Set(status)

	if status && b.OTAIsRunning() {
		b.otaFlash <- struct{}{}
	}
}

func (b *Bind) Broadcast(ctx context.Context, level string, payload interface{}) error {
	return b.MQTT().PublishRaw(ctx, b.config.TopicBroadcast.Format(level), 1, false, payload)
}

func (b *Bind) Restart(ctx context.Context) error {
	return b.MQTT().PublishRaw(ctx, b.config.TopicRestart, 1, false, true)
}

func (b *Bind) Reset(ctx context.Context) error {
	return b.MQTT().Publish(ctx, b.config.TopicReset, true)
}

func (b *Bind) ProtocolVersion() string {
	v, ok := b.DeviceAttribute("homie")
	if ok {
		return v.(string)
	}

	return ""
}

func (b *Bind) ProtocolVersionConstraint(constraint string) bool {
	current, err := version.NewVersion(b.ProtocolVersion())
	if err != nil {
		return false
	}

	constraints, err := version.NewConstraint(constraint)
	if err != nil {
		return false
	}

	return constraints.Check(current)
}

func (b *Bind) bump() {
	b.lastUpdate.Set(time.Now())
}

func (b *Bind) LastUpdate() *time.Time {
	if b.lastUpdate.IsNil() {
		return nil
	}

	return b.lastUpdate.Load()
}

func (b *Bind) nodesList() []*node {
	result := make([]*node, 0)

	b.nodes.Range(func(key, value interface{}) bool {
		result = append(result, value.(*node))
		return true
	})

	return result
}

func (b *Bind) nodesAttributesSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	for _, name := range strings.Split(message.String(), ",") {
		b.nodes.Store(name, &node{
			ID:         atomic.NewStringDefault(name),
			Name:       atomic.NewString(),
			Type:       atomic.NewString(),
			Array:      atomic.NewString(),
			properties: &sync.Map{},
		})

		b.Logger().Debug("Register node", "node", name)
	}

	return nil
}

/*
	homie / device ID / node ID / $node-attribute
	homie / device ID / node ID / property ID
*/
func (b *Bind) nodesSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()

	// skip $fw $stats $implementation etc.
	if strings.HasPrefix(route[2], "$") {
		return nil
	}

	value, ok := b.nodes.Load(route[2])
	if !ok {
		return errors.New("unknown node " + route[2])
	}

	n := value.(*node)

	if strings.HasPrefix(route[3], "$") { // node attribute
		switch route[3] {
		case "$name":
			n.Name.Set(message.String())
		case "$type":
			n.Type.Set(message.String())
		case "$properties":
		case "$array":
		}
	} else {
		property := n.propertyLoadOrStore(route[3])
		property.Value.Set(message.String())
	}

	return nil
}

/*
	homie / device ID / node ID / property ID / $property-attribute
*/
func (b *Bind) nodesPropertySubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	route := message.Topic().Split()

	// skip $fw $stats $implementation etc.
	if strings.HasPrefix(route[2], "$") {
		return nil
	}

	value, ok := b.nodes.Load(route[2])
	if !ok {
		return errors.New("unknown node " + route[2])
	}

	n := value.(*node)
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
