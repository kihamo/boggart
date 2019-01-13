package boggart

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kihamo/boggart/components/mqtt"
)

type BindBase struct {
	mutex sync.RWMutex

	serialNumber string
	status       uint64
}

func (d *BindBase) Init() {
	d.UpdateStatus(BindStatusInitializing)
}

func (d *BindBase) Status() BindStatus {
	return BindStatus(atomic.LoadUint64(&d.status))
}

func (d *BindBase) UpdateStatus(status BindStatus) {
	atomic.StoreUint64(&d.status, uint64(status))
}

func (d *BindBase) SerialNumber() string {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	return d.serialNumber
}

func (d *BindBase) SetSerialNumber(serialNumber string) {
	d.mutex.Lock()
	d.serialNumber = serialNumber
	d.mutex.Unlock()
}

type BindMQTT struct {
	mutex  sync.RWMutex
	client mqtt.Component
}

func (d *BindMQTT) SetMQTTClient(client mqtt.Component) {
	d.mutex.Lock()
	d.client = client
	d.mutex.Unlock()
}

func (d *BindMQTT) MQTTPublish(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) error {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if d.client == nil {
		return errors.New("MQTT client isn't init")
	}

	switch value := payload.(type) {
	case string, []byte:
		// skip
	case float64:
		payload = strconv.FormatFloat(value, 'f', -1, 64)
	case float32:
		payload = strconv.FormatFloat(float64(value), 'f', -1, 64)
	case int64:
		payload = strconv.FormatInt(value, 10)
	case int32:
		payload = strconv.FormatInt(int64(value), 10)
	case int16:
		payload = strconv.FormatInt(int64(value), 10)
	case int8:
		payload = strconv.FormatInt(int64(value), 10)
	case int:
		payload = strconv.FormatInt(int64(value), 10)
	case uint64:
		payload = strconv.FormatUint(value, 10)
	case uint32:
		payload = strconv.FormatUint(uint64(value), 10)
	case uint16:
		payload = strconv.FormatUint(uint64(value), 10)
	case uint8:
		payload = strconv.FormatUint(uint64(value), 10)
	case uint:
		payload = strconv.FormatUint(uint64(value), 10)
	case bool:
		if value {
			payload = []byte(`1`)
		} else {
			payload = []byte(`0`)
		}
	case time.Time:
		payload = value.Format(time.RFC3339)
	case *time.Time:
		payload = value.Format(time.RFC3339)
	default:
		payload = fmt.Sprintf("%s", payload)
	}

	return d.client.Publish(ctx, topic, qos, retained, payload)
}

func (d *BindMQTT) MQTTPublishAsync(ctx context.Context, topic string, qos byte, retained bool, payload interface{}) {
	go func() {
		d.MQTTPublish(ctx, topic, qos, retained, payload)
	}()
}

func CheckSerialNumberInMQTTTopic(bind Bind, topic string, offset int) bool {
	sn := mqtt.NameReplace(bind.SerialNumber())

	if sn == "" {
		return false
	}

	routes := mqtt.RouteSplit(topic)
	if len(routes) < offset {
		return false
	}

	return routes[len(routes)-offset] == sn
}

func WrapMQTTSubscribeDeviceIsOnline(bind Bind, callback mqtt.MessageHandler) mqtt.MessageHandler {
	return func(ctx context.Context, client mqtt.Component, message mqtt.Message) {
		if bind.Status() == BindStatusOnline {
			callback(ctx, client, message)
		}
	}
}
