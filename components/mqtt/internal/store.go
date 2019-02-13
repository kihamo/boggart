package internal

import (
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/kihamo/shadow/components/logging"
)

type Store struct {
	store  mqtt.Store
	logger logging.Logger
}

func NewStore(logger logging.Logger) *Store {
	return &Store{
		store:  mqtt.NewMemoryStore(),
		logger: logger,
	}
}

func (s *Store) Open() {
	s.store.Open()
}

func (s *Store) Put(key string, message packets.ControlPacket) {
	switch m := message.(type) {
	case *packets.ConnackPacket:
		s.logger.Debug("Put connect ack packet")

	case *packets.ConnectPacket:
		s.logger.Debug("Put connect packet")

	case *packets.DisconnectPacket:
		s.logger.Debug("Put disconnect packet")

	case *packets.PingreqPacket:
		s.logger.Debug("Put ping request packet")

	case *packets.PingrespPacket:
		s.logger.Debug("Put ping response packet")

	case *packets.PubackPacket:
		s.logger.Debug("Put publish ack packet",
			"message-id", m.MessageID,
		)

	case *packets.PubcompPacket:
		s.logger.Debug("Put pubcomp packet",
			"message-id", m.MessageID,
		)

	case *packets.PublishPacket:
		s.logger.Debug("Put publish packet",
			"message-id", m.MessageID,
			"topic", m.TopicName,
			"qos", m.Qos,
			"retained", m.Retain,
			"payload", string(m.Payload),
		)

	case *packets.PubrecPacket:
		s.logger.Debug("Put pubrec packet",
			"message-id", m.MessageID,
		)

	case *packets.PubrelPacket:
		s.logger.Debug("Put pubrel packet",
			"message-id", m.MessageID,
		)

	case *packets.SubackPacket:
		s.logger.Debug("Put subscribe ack packet",
			"message-id", m.MessageID,
		)

	case *packets.SubscribePacket:
		s.logger.Debug("Put subscribe packet",
			"message-id", m.MessageID,
			"topics", strings.Join(m.Topics, ";"),
		)

	case *packets.UnsubackPacket:
		s.logger.Debug("Put unsubscribe ack packet",
			"message-id", m.MessageID,
		)

	case *packets.UnsubscribePacket:
		s.logger.Debug("Put unsubscribe packet",
			"message-id", m.MessageID,
			"topic", strings.Join(m.Topics, ";"),
		)

	default:
		s.logger.Debug("Put unknown packet", "key", key, "message", message.String())
	}

	s.store.Put(key, message)
}

func (s *Store) Get(key string) packets.ControlPacket {
	return s.store.Get(key)
}

func (s *Store) All() []string {
	return s.store.All()
}

func (s *Store) Del(key string) {
	s.store.Del(key)
}

func (s *Store) Close() {
	s.store.Close()
}

func (s *Store) Reset() {
	s.store.Reset()
}
