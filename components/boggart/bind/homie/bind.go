package homie

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	a "github.com/kihamo/boggart/components/boggart/atomic"
)

const (
	configNameSeparator       = "."
	configDeviceAttributeName = "implementation.config"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config     *Config
	lastUpdate *a.TimeNull

	deviceAttributes     *sync.Map
	implementationConfig *sync.Map
}

type ImplementationConfig struct {
	Name  string
	Type  string
	Value interface{}
}

func (b *Bind) registerDeviceAttributes(name string, value interface{}) {
	b.deviceAttributes.Store(name, value)

	if name == configDeviceAttributeName {
		md := make(map[string]interface{})
		err := json.Unmarshal([]byte(fmt.Sprintf("%v", value)), &md)
		if err == nil {
			b.configMetadataParse(reflect.ValueOf(md), "")
		}
	}
}

func (b *Bind) DeviceAttribute(key string) (interface{}, bool) {
	return b.deviceAttributes.Load(key)
}

func (b *Bind) DeviceAttributes() map[string]interface{} {
	result := make(map[string]interface{})

	b.deviceAttributes.Range(func(key, value interface{}) bool {
		result[key.(string)] = value
		return true
	})

	return result
}

func (b *Bind) configMetadataParse(e reflect.Value, prefix string) {
	if !e.IsValid() || e.Kind() != reflect.Map || e.Len() == 0 {
		return
	}

	for _, field := range e.MapKeys() {
		value := reflect.ValueOf(e.MapIndex(field).Interface())
		key := prefix + field.String()

		switch value.Kind() {
		case reflect.Map:
			b.configMetadataParse(value, key+configNameSeparator)

		default:
			b.implementationConfig.Store(key, ImplementationConfig{
				Name:  prefix + field.String(),
				Type:  value.Kind().String(),
				Value: value.Interface(),
			})
		}
	}
}

func (b *Bind) ImplementationConfigAll() map[string]ImplementationConfig {
	result := make(map[string]ImplementationConfig)

	b.implementationConfig.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(ImplementationConfig)
		return true
	})

	return result
}

func (b *Bind) ImplementationConfig(key string) (value ImplementationConfig, ok bool) {
	if v, o := b.implementationConfig.Load(key); o {
		return v.(ImplementationConfig), o
	}

	return value, ok
}

func (b *Bind) ImplementationConfigSet(ctx context.Context, key string, value interface{}) (err error) {
	md, ok := b.ImplementationConfig(key)
	if !ok {
		return errors.New("config option " + key + " not found")
	}

	switch md.Type {
	case reflect.Bool.String():
		value, err = strconv.ParseBool(value.(string))

	case reflect.Int.String(), reflect.Int32.String(), reflect.Int64.String():
		value, err = strconv.ParseInt(value.(string), 10, 64)

	case reflect.Uint.String(), reflect.Uint32.String(), reflect.Uint64.String():
		value, err = strconv.ParseUint(value.(string), 10, 64)

	case reflect.Float32.String(), reflect.Float64.String():
		value, err = strconv.ParseFloat(value.(string), 64)

	default:
		value = value.(string)
	}

	if err != nil {
		return err
	}

	payload := make(map[string]interface{}, 1)
	levels := strings.Split(key, configNameSeparator)

	for i := len(levels) - 1; i >= 0; i-- {
		if i == len(levels)-1 {
			payload[levels[i]] = value
		} else {
			payload = map[string]interface{}{
				levels[i]: payload,
			}
		}
	}

	pl, err := json.Marshal(&payload)
	if err != nil {
		return err
	}

	return b.MQTTPublish(ctx, MQTTPublishTopicConfigSet.Format(b.config.BaseTopic, b.SerialNumber()), pl)
}

func (b *Bind) Broadcast(ctx context.Context, level string, payload interface{}) error {
	return b.MQTTPublishRaw(ctx, MQTTPublishTopicBroadcast.Format(b.config.BaseTopic, level), 1, false, payload)
}

func (b *Bind) Restart(ctx context.Context) error {
	return b.MQTTPublishRaw(ctx, MQTTPublishTopicRestart.Format(b.config.BaseTopic, b.SerialNumber()), 1, false, true)
}

func (b *Bind) Reset(ctx context.Context) error {
	return b.MQTTPublish(ctx, MQTTPublishTopicReset.Format(b.config.BaseTopic, b.SerialNumber()), true)
}

func (b *Bind) bump() {
	b.lastUpdate.Set(time.Now())
}

func (b *Bind) LastUpdate() *time.Time {
	if b.lastUpdate.IsNil() {
		return nil
	}

	t := b.lastUpdate.Load()
	return &t
}
