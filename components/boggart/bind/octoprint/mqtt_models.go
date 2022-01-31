package octoprint

import (
	"encoding/json"
	"time"
)

type Payload struct {
	Timestamp          time.Time
	data               map[string]interface{}
	timestampFieldName string
}

func (p *Payload) Time() time.Time {
	return p.Timestamp
}

func (p *Payload) SetTimestampFieldName(name string) {
	p.timestampFieldName = name
}

func (p *Payload) UnmarshalJSON(b []byte) (err error) {
	if err = json.Unmarshal(b, &p.data); err != nil {
		return err
	}

	fieldName := "_timestamp"
	if p.timestampFieldName != "" {
		fieldName = p.timestampFieldName
	}

	if raw, ok := p.data[fieldName]; ok {
		if val, ok := raw.(float64); ok {
			p.Timestamp = time.Unix(int64(val), 0)
		}
	}

	return nil
}

type Temperature struct {
	*Payload

	Actual float64 `json:"actual"`
	Target float64 `json:"target"`
}

func NewTemperature(name string) *Temperature {
	return &Temperature{
		Payload: &Payload{
			timestampFieldName: name,
		},
	}
}

func (t *Temperature) UnmarshalJSON(b []byte) error {
	if t.Payload == nil {
		t.Payload = &Payload{}
	}

	if err := t.Payload.UnmarshalJSON(b); err != nil {
		return err
	}

	if value, ok := t.Payload.data["actual"]; ok {
		t.Actual = value.(float64)
	}

	if value, ok := t.Payload.data["target"]; ok {
		t.Target = value.(float64)
	}

	return nil
}

type Event struct {
	*Payload

	Event string `json:"_event"`
}

func (e *Event) UnmarshalJSON(b []byte) error {
	if e.Payload == nil {
		e.Payload = &Payload{}
	}

	if err := e.Payload.UnmarshalJSON(b); err != nil {
		return err
	}

	if value, ok := e.Payload.data["_event"]; ok {
		e.Event = value.(string)
	}

	return nil
}

type EventSettingsUpdated struct {
	*Event

	ConfigHash    string `json:"config_hash"`
	EffectiveHash string `json:"effective_hash"`
}

func (e *EventSettingsUpdated) UnmarshalJSON(b []byte) error {
	if e.Event == nil {
		e.Event = &Event{}
	}

	if err := e.Event.UnmarshalJSON(b); err != nil {
		return err
	}

	if value, ok := e.Payload.data["config_hash"]; ok {
		e.ConfigHash = value.(string)
	}

	if value, ok := e.Payload.data["effective_hash"]; ok {
		e.EffectiveHash = value.(string)
	}

	return nil
}
