package openhab

import (
	"github.com/pborman/uuid"
)

type thing struct {
	id         string
	bindingID  string
	typeID     string
	label      string
	bridge     string
	location   string
	parameters *Parameters
}

func newThing(id, bindingID, typeID string) *thing {
	return (&thing{
		bindingID:  IDReplace(bindingID),
		typeID:     IDReplace(typeID),
		parameters: newParameters(),
	}).
		withID(id)
}

func (t *thing) ThingID() string {
	if t == nil {
		return ""
	}

	return t.bindingID + ":" + t.typeID + ":" + t.id
}

func (t *thing) withID(id string) *thing {
	t.id = IDReplace(id)
	return t
}

func (t *thing) withLabel(label string) *thing {
	t.label = label
	return t
}

func (t *thing) withBridge(bridge string) *thing {
	t.bridge = bridge
	return t
}

func (t *thing) withLocation(location string) *thing {
	t.location = location
	return t
}

// Thing <binding_id>:<type_id>:<thing_id> "Label" @ "Location" [ <parameters> ]
func (t *thing) String() string {
	if t.id == "" {
		t.withID(uuid.New())
	}

	s := t.ThingID()

	if t.label != "" {
		s += " \"" + t.label + "\""
	}

	if t.bridge != "" {
		s += " (" + t.bridge + ")"
	}

	if t.location != "" {
		s += " @ \"" + t.location + "\""
	}

	if t.parameters != nil {
		if params := t.parameters.String(); params != "" {
			s += " [" + params + "]"
		}
	}

	return s
}
