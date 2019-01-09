package devices

import (
	"github.com/kihamo/boggart/components/boggart"
)

func init() {
	boggart.RegisterKind("camera_hikvision", CameraHikVision{})
	boggart.RegisterKind("led_wifi", WiFiLED{})
	boggart.RegisterKind("internet_provider_softvideo", SoftVideoInternet{})
	boggart.RegisterKind("remote_control_broadlink_rm", BroadlinkRMRemoteControl{})
	boggart.RegisterKind("router_mikrotik", MikrotikRouter{})
	boggart.RegisterKind("sensor_ds18b20", DS18B20Sensor{})
	boggart.RegisterKind("smart_speaker_google_home_mini", GoogleHomeMiniSmartSpeaker{})
	boggart.RegisterKind("socket_broadlink_sp3s", BroadlinkSP3SSocket{})
	boggart.RegisterKind("tv_lg_webos", LGWebOSTV{})
	boggart.RegisterKind("tv_samsung", SamsungTV{})
	boggart.RegisterKind("ups_nut", UPSNUT{})
}
