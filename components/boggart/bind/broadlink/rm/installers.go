package rm

import (
	"context"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
		installer.SystemDevice,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, s installer.System) ([]installer.Step, error) {
	cfg := b.config()

	if s == installer.SystemDevice {
		return []installer.Step{{
			FilePath:    "setup.go",
			Description: "Setup script",
			Content: `package main

import (
	"fmt"
	"github.com/kihamo/boggart/providers/broadlink"
	"log"
)

func main() {
	err := broadlink.SetupWiFi("wifi_sid", "wifi_password", broadlink.WifiSecurityWPA2)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Setup success")
}`,
		}, {
			Content: `1. Долго зажать кнопку на устройстве, светодиод должен начать работать в режиме частое моргание - пауза - частое мограние
2. Подключиться к незащищенной точке BroadlinkProv
3. Запустить go run setup.go
4. Заменить wifi_sid и wifi_password на соответствующие значения
4. Подключиться к сети указанной в setup.go и через Discover обнаружить устройство (должно быть по адресу ` + cfg.Host + `)`,
		}}, nil
	}

	meta := b.Meta()
	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	id := meta.ID()

	channels := make([]*openhab.Channel, 0)

	if _, ok := b.provider.(SupportCapture); ok {
		const (
			idCaptureSwitch     = "CaptureSwitch"
			idCaptureResultType = "CaptureResult_Type"
			idCaptureResultCode = "CaptureResult_Code"
		)

		channels = append(channels,
			openhab.NewChannel(idCaptureSwitch, openhab.ChannelTypeSwitch).
				WithStateTopic(cfg.TopicCaptureState.Format(id)).
				WithCommandTopic(cfg.TopicCaptureSwitch.Format(id)).
				WithOn("true").
				WithOff("false").
				AddItems(
					openhab.NewItem(itemPrefix+idCaptureSwitch, openhab.ItemTypeSwitch).
						WithLabel("Capture switch").
						WithIcon("text"),
				),
			openhab.NewChannel(idCaptureResultType, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicCaptureResult.Format(id)).
				WithTransformationPattern("JSONPATH:$.type").
				AddItems(
					openhab.NewItem(itemPrefix+idCaptureResultType, openhab.ItemTypeString).
						WithLabel("Capture type").
						WithIcon("text"),
				),
			openhab.NewChannel(idCaptureResultCode, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicCaptureResult.Format(id)).
				WithTransformationPattern("JSONPATH:$.code").
				AddItems(
					openhab.NewItem(itemPrefix+idCaptureResultCode, openhab.ItemTypeString).
						WithLabel("Capture code").
						WithIcon("text"),
				),
		)
	}

	const (
		transformSend installer.Path = openhab.DirectoryTransform + "broadlink_rm_send.js"
	)

	var sendExists bool

	if _, ok := b.provider.(SupportIR); ok {
		const idSendIR = "Send_IR"
		sendExists = true

		channels = append(channels,
			openhab.NewChannel(idSendIR, openhab.ChannelTypeString).
				WithCommandTopic(cfg.TopicIR.Format(id)).
				WithTransformationPatternOut("JS:"+transformSend.Base()+"?set_count=0").
				AddItems(
					openhab.NewItem(itemPrefix+idSendIR, openhab.ItemTypeString).
						WithLabel("IR command").
						WithIcon("text"),
				),
		)
	}

	if _, ok := b.provider.(SupportRF315Mhz); ok {
		const idSendRF315 = "Send_RF315"
		sendExists = true

		channels = append(channels,
			openhab.NewChannel(idSendRF315, openhab.ChannelTypeString).
				WithCommandTopic(cfg.TopicRF315.Format(id)).
				WithTransformationPatternOut("JS:"+transformSend.Base()+"?set_count=0").
				AddItems(
					openhab.NewItem(itemPrefix+idSendRF315, openhab.ItemTypeString).
						WithLabel("RF315 command").
						WithIcon("text"),
				),
		)
	}

	if _, ok := b.provider.(SupportRF433Mhz); ok {
		const idSendRF433 = "Send_RF433"
		sendExists = true

		channels = append(channels,
			openhab.NewChannel(idSendRF433, openhab.ChannelTypeString).
				WithCommandTopic(cfg.TopicRF433.Format(id)).
				WithTransformationPatternOut("JS:"+transformSend.Base()+"?set_count=0").
				AddItems(
					openhab.NewItem(itemPrefix+idSendRF433, openhab.ItemTypeString).
						WithLabel("RF433 command").
						WithIcon("text"),
				),
		)
	}

	steps := make([]installer.Step, 0, 1)
	if sendExists {
		steps = append(steps, installer.Step{
			FilePath: transformSend,
			Content: `(function(value, count) {
    return JSON.stringify({ code: value, count: parseInt(count, 10) });
})(input, set_count)`,
		})
	}

	return openhab.StepsByBind(b, steps, channels...)
}
