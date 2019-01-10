package devices

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterDeviceType("camera_hikvision", CameraHikVision{})
	boggart.RegisterDeviceType("led_wifi", WiFiLED{})
	boggart.RegisterDeviceType("softvideo", SoftVideo{})
	boggart.RegisterDeviceType("broadlink_rm", BroadlinkRM{})
	boggart.RegisterDeviceType("mikrotik", Mikrotik{})
	boggart.RegisterDeviceType("sensor_ds18b20", DS18B20Sensor{})
	boggart.RegisterDeviceType("smart_speaker_google_home_mini", GoogleHomeMiniSmartSpeaker{})
	boggart.RegisterDeviceType("broadlink_sp3s", BroadlinkSP3S{})
	boggart.RegisterDeviceType("tv_lg_webos", LGWebOSTV{})
	boggart.RegisterDeviceType("tv_samsung", SamsungTV{})
	boggart.RegisterDeviceType("nut", NUT{})
}
