package grafana

import (
	"github.com/kihamo/boggart/components/boggart/config_generators"
	"github.com/kihamo/boggart/components/boggart/config_generators/openhab"
)

func (b *Bind) GenerateConfigOpenHab() ([]generators.Step, error) {
	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())

	const idAnnotation = "Annotation"

	return openhab.StepsByBind(b, []generators.Step{{
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
