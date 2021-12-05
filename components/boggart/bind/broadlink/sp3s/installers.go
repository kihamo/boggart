package sp3s

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

	itemPrefix := openhab.ItemPrefixFromBindMeta(b.Meta())
	mac := cfg.MAC.String()

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanWatts),
	},
		openhab.NewChannel("Status", openhab.ChannelTypeSwitch).
			WithStateTopic(cfg.TopicState.Format(mac)).
			WithCommandTopic(cfg.TopicSet.Format(mac)).
			WithOn("true").
			WithOff("false").
			AddItems(
				openhab.NewItem(itemPrefix+"Status", openhab.ItemTypeSwitch).
					WithLabel("Status []"),
			),
		openhab.NewChannel("Power", openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicPower.Format(mac)).
			AddItems(
				openhab.NewItem(itemPrefix+"Power", openhab.ItemTypeNumber).
					WithLabel("Power [JS("+openhab.StepDefaultTransformHumanWatts.Base()+"):%s]").
					WithIcon("energy"),
			),
	)
}
