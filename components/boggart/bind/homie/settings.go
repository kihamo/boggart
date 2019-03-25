package homie

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/kihamo/boggart/components/mqtt"
)

const (
	settingsMQTTPublishTopicGet = MQTTPrefixImpl + "config"
	settingsMQTTPublishTopicSet = MQTTPrefixImpl + "config/set"
)

type SettingsOption struct {
	Name  string
	Type  string
	Value interface{}
}

func (b *Bind) settingsParse(e reflect.Value, prefix string) {
	if !e.IsValid() || e.Kind() != reflect.Map || e.Len() == 0 {
		return
	}

	for _, field := range e.MapKeys() {
		value := reflect.ValueOf(e.MapIndex(field).Interface())
		key := prefix + field.String()

		switch value.Kind() {
		case reflect.Map:
			b.settingsParse(value, key+configNameSeparator)

		default:
			b.settings.Store(key, SettingsOption{
				Name:  prefix + field.String(),
				Type:  value.Kind().String(),
				Value: value.Interface(),
			})
		}
	}
}

func (b *Bind) SettingsAll() map[string]SettingsOption {
	result := make(map[string]SettingsOption)

	b.settings.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(SettingsOption)
		return true
	})

	return result
}

func (b *Bind) SettingsGet(key string) (value SettingsOption, ok bool) {
	if v, o := b.settings.Load(key); o {
		return v.(SettingsOption), o
	}

	return value, ok
}

func (b *Bind) SettingsSet(ctx context.Context, key string, value interface{}) (err error) {
	md, ok := b.SettingsGet(key)
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

	return b.MQTTPublish(ctx, settingsMQTTPublishTopicSet.Format(b.config.BaseTopic, b.SerialNumber()), pl)
}

func (b *Bind) settingsSubscriber(_ context.Context, _ mqtt.Component, message mqtt.Message) error {
	md := make(map[string]interface{})
	err := message.UnmarshalJSON(&md)
	if err == nil {
		b.settingsParse(reflect.ValueOf(md), "")
	}

	return err
}
