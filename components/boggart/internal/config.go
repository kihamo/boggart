package internal

import (
	"time"

	"github.com/davecheney/gpio"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/protocols/rs485"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			boggart.ConfigRS485Address,
			config.ValueTypeString,
			rs485.DefaultSerialAddress,
			"Serial port address for RS485",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigRS485Timeout,
			config.ValueTypeDuration,
			rs485.DefaultTimeout,
			"Serial port timeout for RS485",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigDoorsEnabled,
			config.ValueTypeBool,
			false,
			"Enabled doors provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigDoorsEntrancePin,
			config.ValueTypeInt,
			gpio.GPIO17,
			"Pin for entrance door reed switch",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMercuryEnabled,
			config.ValueTypeBool,
			false,
			"Enabled mercury provider",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMercuryRepeatInterval,
			config.ValueTypeDuration,
			time.Minute*2,
			"Repeat interval for mercury provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMercuryDeviceAddress,
			config.ValueTypeString,
			nil,
			"Mercury device address in format XXXXXX (last 6 digits of device serial number)",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMikrotikEnabled,
			config.ValueTypeBool,
			false,
			"Enabled Mikrotik provider",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMikrotikRepeatInterval,
			config.ValueTypeDuration,
			time.Minute*5,
			"Repeat interval for Mikrotik provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMikrotikAddress,
			config.ValueTypeString,
			"192.168.88.1:8728",
			"Mikrotik API address in format host:port",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMikrotikUsername,
			config.ValueTypeString,
			"admin",
			"Username of Mikrotik API account",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMikrotikPassword,
			config.ValueTypeString,
			nil,
			"Password of Mikrotik API account",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			boggart.ConfigMikrotikTimeout,
			config.ValueTypeDuration,
			time.Second*10,
			"Request timeout",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMobileEnabled,
			config.ValueTypeBool,
			false,
			"Enabled Megafon mobile provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMobileRepeatInterval,
			config.ValueTypeDuration,
			time.Minute*30,
			"Repeat interval for mobile provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMobileMegafonPhone,
			config.ValueTypeString,
			nil,
			"Phone of Megafon account",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigMobileMegafonPassword,
			config.ValueTypeString,
			nil,
			"Password of Megafon account",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			boggart.ConfigPulsarEnabled,
			config.ValueTypeBool,
			false,
			"Enabled pulsar provider",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarRepeatInterval,
			config.ValueTypeDuration,
			time.Minute*15,
			"Repeat interval for pulsar provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarHeatMeterAddress,
			config.ValueTypeString,
			nil,
			"Pulsar address HEX value (AABBCCDD). If empty system try to find device",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarColdWaterPulseInput,
			config.ValueTypeInt64,
			pulsar.Input1,
			"Pulsar input number of cold water",
			true,
			[]string{config.ViewEnum},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{pulsar.Input1, "#1"},
					{pulsar.Input2, "#2"},
				},
			}),
		config.NewVariable(
			boggart.ConfigPulsarColdWaterStartValue,
			config.ValueTypeFloat64,
			0,
			"Pulsar start value of cold water (in m3)",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarHotWaterPulseInput,
			config.ValueTypeInt64,
			pulsar.Input2,
			"Pulsar input number of hot water",
			true,
			[]string{config.ViewEnum},
			map[string]interface{}{
				config.ViewOptionEnumOptions: [][]interface{}{
					{pulsar.Input1, "#1"},
					{pulsar.Input2, "#2"},
				},
			}),
		config.NewVariable(
			boggart.ConfigPulsarHotWaterStartValue,
			config.ValueTypeFloat64,
			0,
			"Pulsar start value of hot water (in m3)",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoEnabled,
			config.ValueTypeBool,
			false,
			"Enabled SoftVideo provider",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoRepeatInterval,
			config.ValueTypeDuration,
			time.Hour*12,
			"Repeat interval for SoftVideo provider",
			false,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoLogin,
			config.ValueTypeString,
			nil,
			"Login of SoftVideo account",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigSoftVideoPassword,
			config.ValueTypeString,
			nil,
			"Password of SoftVideo account",
			true,
			[]string{config.ViewPassword},
			nil),
		config.NewVariable(
			boggart.ConfigMonitoringExternalURL,
			config.ValueTypeString,
			nil,
			"Monitoring external URL",
			false,
			nil,
			nil),
	}
}

func (c *Component) GetConfigWatchers() []config.Watcher {
	return []config.Watcher{
		config.NewWatcher(c.GetName(), []string{
			boggart.ConfigRS485Timeout,
			boggart.ConfigRS485Address,
		}, c.watchConnectionRS485),
	}
}

func (c *Component) watchConnectionRS485(_ string, _ interface{}, _ interface{}) {
	c.initConnectionRS485()
}
