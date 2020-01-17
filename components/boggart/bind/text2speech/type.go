package text2speech

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/protocols/swagger"
	"github.com/kihamo/boggart/providers/yandex_speechkit_cloud"
)

type Type struct {
	boggart.BindTypeWidget
}

func (t Type) CreateBind(c interface{}) (boggart.Bind, error) {
	cfg := c.(*Config)

	cfg.TopicGenerateBinaryOptions = cfg.TopicGenerateBinaryOptions.Format(cfg.Key)
	cfg.TopicGenerateBinaryText = cfg.TopicGenerateBinaryText.Format(cfg.Key)
	cfg.TopicGenerateURLOptions = cfg.TopicGenerateURLOptions.Format(cfg.Key)
	cfg.TopicGenerateURLText = cfg.TopicGenerateURLText.Format(cfg.Key)
	cfg.TopicURL = cfg.TopicURL.Format(cfg.Key)
	cfg.TopicBinary = cfg.TopicBinary.Format(cfg.Key)

	bind := &Bind{
		config: cfg,
	}

	l := swagger.NewLogger(
		func(message string) {
			bind.Logger().Info(message)
		},
		func(message string) {
			bind.Logger().Debug(message)
		})

	bind.provider = yandex_speechkit_cloud.New(cfg.Debug, l)

	return bind, nil
}
