package syslog

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/kihamo/boggart/components/boggart/di"
	"gopkg.in/mcuadros/go-syslog.v2"
	"gopkg.in/mcuadros/go-syslog.v2/format"
)

type Bind struct {
	status uint32 // 0 - default, 1 - started, 2 - failed

	di.ConfigBind
	di.LoggerBind
	di.ProbesBind
	di.MQTTBind

	config *Config
	mutex  sync.Mutex
	server *syslog.Server
	addr   string
}

func (b *Bind) Run() error {
	atomic.StoreUint32(&b.status, 0)

	go func() {
		atomic.StoreUint32(&b.status, 1)
		defer atomic.StoreUint32(&b.status, 2)

		server := syslog.NewServer()
		server.SetFormat(syslog.Automatic)
		server.SetTimeout(int64(b.config.Timeout.Seconds()))
		server.SetHandler(b)

		if err := server.ListenUDP(b.addr); err != nil {
			b.Logger().Error("Failed listen with error " + err.Error())
			return
		}

		if err := server.Boot(); err != nil {
			b.Logger().Error("Failed boot with error " + err.Error())
			return
		}

		b.mutex.Lock()
		b.server = server
		b.mutex.Unlock()

		server.Wait()
	}()

	return nil
}

func (b *Bind) Close() error {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if b.server != nil {
		return b.server.Kill()
	}

	return nil
}

func (b *Bind) Handle(message format.LogParts, length int64, err error) {
	var tag string
	if val, ok := message["tag"]; ok {
		tag = fmt.Sprintf("%v", val)
	}

	metricHandledMessage.With("address", b.addr, "tag", tag).Inc()

	if err != nil {
		b.Logger().Error("Handler with error",
			"error", err,
			"length", length,
			"message", message,
		)
	} else {
		b.Logger().Debug("Receive message",
			"length", length,
			"message", message,
		)
	}

	var fields map[string]interface{}

	if len(b.config.Tags) > 0 {
		fields = make(map[string]interface{}, len(b.config.Tags))

		for _, tag := range b.config.Tags {
			if value, ok := message[tag]; ok {
				fields[tag] = value
			}
		}
	} else {
		fields = message
	}

	topic := b.config.Topic
	if topic == "" {
		if val, ok := message["hostname"]; ok {
			topic = fmt.Sprintf("%v", val)
		}
	}

	if topic == "" {
		topic = b.addr
	}

	payload, err := json.Marshal(fields)
	if err != nil {
		b.Logger().Error("JSON message failed",
			"error", err,
			"message", message,
		)

		return
	}

	_ = b.MQTT().PublishAsyncWithoutCache(context.Background(), b.config.TopicMessage.Format(topic), payload)
}
