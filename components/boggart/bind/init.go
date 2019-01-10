package bind

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("broadlink_rm", BroadlinkRM{})
	boggart.RegisterDeviceType("broadlink_sp3s", BroadlinkSP3S{})
	boggart.RegisterDeviceType("ds18b20", DS18B20{})
	boggart.RegisterDeviceType("google_home_mini", GoogleHomeMini{})
	boggart.RegisterDeviceType("gpio", GPIO{})
	boggart.RegisterDeviceType("hikvision", HikVision{})
	boggart.RegisterDeviceType("led_wifi", WiFiLED{})
	boggart.RegisterDeviceType("lg_webos", LGWebOS{})
	boggart.RegisterDeviceType("mikrotik", Mikrotik{})
	boggart.RegisterDeviceType("nut", NUT{})
	boggart.RegisterDeviceType("pulsar_heat_meter", PulsarHeatMeter{})
	boggart.RegisterDeviceType("samsung_tizen", SamsungTizen{})
	boggart.RegisterDeviceType("softvideo", SoftVideo{})

}
