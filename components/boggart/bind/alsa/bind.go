package alsa

import (
	"time"

	"github.com/kihamo/boggart/components/boggart"
	a "github.com/kihamo/boggart/components/voice/players/alsa"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	player *a.Player

	status          int64
	volume          int64
	mute            int64
	updaterInterval time.Duration
}
