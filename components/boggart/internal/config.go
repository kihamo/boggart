package internal

import (
	"time"

	"github.com/davecheney/gpio"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/pulsar"
	"github.com/kihamo/shadow/components/config"
)

func (c *Component) GetConfigVariables() []config.Variable {
	return []config.Variable{
		config.NewVariable(
			boggart.ConfigPulsarEnabled,
			config.ValueTypeBool,
			true,
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
			boggart.ConfigPulsarSerialAddress,
			config.ValueTypeString,
			pulsar.DefaultSerialAddress,
			"Serial port address",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarSerialTimeout,
			config.ValueTypeDuration,
			pulsar.DefaultTimeout,
			"Serial port address",
			true,
			nil,
			nil),
		config.NewVariable(
			boggart.ConfigPulsarDeviceAddress,
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
			true,
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
			boggart.ConfigMikrotikEnabled,
			config.ValueTypeBool,
			true,
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
			boggart.ConfigDoorsEnabled,
			config.ValueTypeBool,
			true,
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
	}
}
