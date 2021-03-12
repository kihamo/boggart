package openhab

import (
	"errors"
	"strconv"
	"strings"

	"github.com/eclipse/paho.mqtt.golang"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/components/boggart/installer"
)

const (
	BindingID = "mqtt"

	DirectoryThings    = "things/"
	DirectoryItems     = "items/"
	DirectoryTransform = "transform/"

	StepDefaultTransformHumanBytes   = DirectoryTransform + "human_bytes.js"
	StepDefaultTransformHumanWatts   = DirectoryTransform + "human_watts.js"
	StepDefaultTransformHumanSeconds = DirectoryTransform + "human_seconds.js"
)

var (
	replacerID = strings.NewReplacer(
		":", "_",
		"{", "_",
		"}", "_",
		"[", "_",
		"]", "_",
		"@", "_",
		" ", "_",
		"\"", "_",
		"-", "_",
	)

	defaultSteps = map[string]installer.Step{
		StepDefaultTransformHumanBytes: {
			FilePath: StepDefaultTransformHumanBytes,
			Content: `(function(i) {
    var
        d = 2,
        e = ['bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'],
        c = 1024;

    if (0 === i || '0' === i) {
        return '0 ' + e[0];
    }

    var f = Math.floor(Math.log(i) / Math.log(c));

    return parseFloat((i / Math.pow(c, f)).toFixed(d)) + ' ' + e[f];
})(input);`,
		},
		StepDefaultTransformHumanWatts: {
			FilePath: StepDefaultTransformHumanWatts,
			Content: `(function(i) {
    var
        d = 2,
        e = ['watts', 'KW', 'MW', 'GW', 'TW', 'PW', 'EW', 'ZW', 'YW'],
        c = 1e3;

    if (0 === i || '0' === i) {
        return '0 ' + e[0];
    }

    var f = Math.floor(Math.log(i) / Math.log(c));

    return parseFloat((i / Math.pow(c, f)).toFixed(d)) + ' ' + e[f];
})(input);`,
		},
		StepDefaultTransformHumanSeconds: {
			FilePath: StepDefaultTransformHumanSeconds,
			Content: `(function(i) {
    var val = parseInt(i);
    var days = 0;
    var hours = 0;
    var minutes = 0;
    var seconds = 0;

    if (val >= 86400) {
        days = Math.floor(val / 86400);
        val = val - (days * 86400);
    }
    if (val >= 3600) {
        hours = Math.floor(val / 3600);
        val = val - (hours * 3600);
    }
    if (val >= 60) {
        minutes = Math.floor(val / 60);
        val = val - (minutes * 60);
    }

    seconds = Math.floor(val);

    var stringDays = '';
    var stringHours = '';
    var stringMinutes = '';
    var stringSeconds = '';

    if (days === 1) {
        stringDays = '1 day ';
    } else if (days > 1) {
        stringDays = days + ' days ';
    }

    if (hours === 1) {
        stringHours = '1 hour ';
    } else if (hours > 1) {
        stringHours = hours + ' hours ';
    }

    if (minutes === 1) {
        stringMinutes = '1 minute ';
    } else if (minutes > 1) {
        stringMinutes = minutes + ' minutes ';
    }

    if (seconds === 1) {
        stringSeconds = '1 second';
    } else if (seconds > 1) {
        stringSeconds = seconds + ' seconds';
    }

    var returnString =  stringDays + stringHours + stringMinutes + stringSeconds;
    return returnString.trim();

})(input);`,
		},
		// TODO: human duration
	}
)

func IDNormalize(id string) string {
	return replacerID.Replace(id)
}

func IDNormalizeCamelCase(id string) string {
	id = IDNormalize(id)
	id = strings.ReplaceAll(id, "_", " ")
	id = strings.Title(id)
	id = strings.ReplaceAll(id, " ", "")

	return id
}

func BrokerFromClientOptionsReader(ops *mqtt.ClientOptionsReader) *Broker {
	server := ops.Servers()[0]
	port, _ := strconv.Atoi(server.Port())
	tsl := ops.TLSConfig()

	return NewBroker(ops.ClientID(), server.Hostname()).
		WithLabel("Auto generate from boggart").
		// WithClientID("openhab").
		WithKeepAlive(int(ops.KeepAlive().Seconds())).
		WithUsername(ops.Username()).
		WithPassword(ops.Password()).
		WithTimeoutInMs(ops.WriteTimeout().Milliseconds()).
		WithPort(port).
		WithSecure(tsl != nil)
}

func StepsByBind(bind boggart.Bind, steps []installer.Step, channels ...*Channel) ([]installer.Step, error) {
	ctrMQTT, ok := di.MQTTContainerBind(bind)
	if !ok {
		return nil, errors.New("bind not supported MQTT container")
	}

	ctrMeta, ok := di.MetaContainerBind(bind)
	if !ok {
		return nil, errors.New("bind not supported Meta container")
	}

	opts, err := ctrMQTT.ClientOptions()
	if err != nil {
		return nil, err
	}

	broker := BrokerFromClientOptionsReader(opts)

	thing := GenericThingFromBindMeta(ctrMeta).
		WithBroker(broker).
		AddChannels(
			BindStatusChannel(ctrMeta),
			BindReloadChannel(ctrMeta),
		)

	if ctrMeta.SerialNumber() != "" {
		thing.AddChannels(BindSerialNumberChannel(ctrMeta))
	}

	if ctrMeta.MAC() != nil {
		thing.AddChannels(BindMACChannel(ctrMeta))
	}

	filePrefix := FilePrefixFromBindMeta(ctrMeta)

	thing.AddChannels(channels...)

	list := make([]installer.Step, 0, len(steps)+3)
	list = append(list, steps...)

	if content := broker.String(); content != "" {
		list = append(list, installer.Step{
			FilePath: DirectoryThings + "broker.things",
			Content:  content,
		})
	}

	if content := thing.String(); content != "" {
		list = append(list, installer.Step{
			FilePath: DirectoryThings + filePrefix + ".things",
			Content:  content,
		})
	}

	if content := thing.Items().String(); content != "" {
		list = append(list, installer.Step{
			FilePath: DirectoryItems + filePrefix + ".items",
			Content:  content,
		})
	}

	return list, nil
}

func BindStatusChannel(meta *di.MetaContainer) *Channel {
	const id = "BindStatus"

	return NewChannel(id, ChannelTypeString).
		WithStateTopic(meta.MQTTTopicStatus()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+id, ItemTypeString).
				WithLabel("Bind status [%s]").
				WithIcon("text"),
		)
}

func BindReloadChannel(meta *di.MetaContainer) *Channel {
	const id = "BindReload"

	return NewChannel(id, ChannelTypeSwitch).
		WithStateTopic(meta.MQTTTopicReload()).
		WithCommandTopic(meta.MQTTTopicReload()).
		WithOn("reload").
		WithOff("done").
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+id, ItemTypeSwitch).
				WithLabel("Bind reload [%s]").
				WithIcon("switch"),
		)
}

func BindSerialNumberChannel(meta *di.MetaContainer) *Channel {
	const id = "BindSerialNumber"

	return NewChannel(id, ChannelTypeString).
		WithStateTopic(meta.MQTTTopicSerialNumber()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+id, ItemTypeString).
				WithLabel("Bind serial number [%s]").
				WithIcon("text"),
		)
}

func BindMACChannel(meta *di.MetaContainer) *Channel {
	return NewChannel("BindMAC", ChannelTypeString).
		WithStateTopic(meta.MQTTTopicMAC()).
		AddItems(
			NewItem(ItemPrefixFromBindMeta(meta)+"BindMAC", ItemTypeString).
				WithLabel("Bind MAC address [%s]").
				WithIcon("text"),
		)
}

func GenericThingFromBindMeta(meta *di.MetaContainer) *GenericThing {
	return NewGenericThing(strings.ToLower(meta.ID())).
		WithLabel(meta.Description())
}

func ItemPrefixFromBindMeta(meta *di.MetaContainer) string {
	return IDNormalizeCamelCase(meta.ID()) + "_"
}

func FilePrefixFromBindMeta(meta *di.MetaContainer) string {
	return IDNormalize("bind_" + strings.ToLower(meta.Type()))
}

func StepDefault(name string) installer.Step {
	return defaultSteps[name]
}
