package voice

import (
	yandex "github.com/kihamo/boggart/components/voice/providers/yandex_speechkit_cloud"
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	TextToSpeechProvider() *yandex.YandexSpeechKitCloud
}
