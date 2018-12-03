package mqtt

type HasSubscribers interface {
	MQTTSubscribers() []Subscriber
}

type Subscriber interface {
	Topic() string
	QOS() byte
	Callback() MessageHandler
}

type SubscriberSimple struct {
	topic    string
	qos      byte
	callback MessageHandler
}

func NewSubscriber(topic string, qos byte, callback MessageHandler) *SubscriberSimple {
	return &SubscriberSimple{
		topic:    topic,
		qos:      qos,
		callback: callback,
	}
}

func (s *SubscriberSimple) Topic() string {
	return s.topic
}

func (s *SubscriberSimple) QOS() byte {
	return s.qos
}

func (s *SubscriberSimple) Callback() MessageHandler {
	return s.callback
}
