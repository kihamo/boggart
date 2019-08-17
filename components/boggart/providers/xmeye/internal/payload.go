package internal

import (
	"bytes"
	"encoding/json"
)

type Payload struct {
	*bytes.Buffer
}

func NewPayload() *Payload {
	return &Payload{
		Buffer: bytes.NewBuffer(nil),
	}
}

func (p *Payload) UnmarshalJSON(v interface{}) error {
	// обрезаем признак конца строки
	payload := p.Bytes()
	payload = payload[:len(payload)-len(payloadEOF)]

	return json.Unmarshal(payload, v)
}
