package voice

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component
	Speaker
}
