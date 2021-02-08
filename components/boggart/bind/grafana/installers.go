package grafana

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(context.Context, installer.System) ([]installer.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	const idAnnotation = "Annotation"

	return openhab.StepsByBind(b, []installer.Step{{
		FilePath: "grafana_annotation.js",
		Content: `(function(i) {
    try {
        var annotation = JSON.parse(i);

        if (typeof annotation.tags !== 'undefined') {
            annotation.tags.push('openhab');
        } else {
            annotation.tags = ['openhab'];
        }

        return JSON.stringify(annotation);
    } catch(e) {
    }

    return '{"text":"' + i + '","tags":["openhab"]}';
})(input);`,
	}},
		openhab.NewChannel(idAnnotation, openhab.ChannelTypeString).
			WithCommandTopic(b.config.TopicAnnotation).
			WithTransformationPatternOut("JS:grafana_annotation.js").
			AddItems(
				openhab.NewItem(itemPrefix+idAnnotation, openhab.ItemTypeString).
					WithLabel("Annotation").
					WithIcon("text"),
			),
	)
}
