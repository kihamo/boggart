package apcupsd

import (
	"errors"

	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/apcupsd"
	"github.com/kihamo/boggart/protocols/apcupsd/file"
	"github.com/kihamo/boggart/protocols/apcupsd/nis"
)

type Bind struct {
	di.ConfigBind
	di.LoggerBind
	di.MetaBind
	di.MetricsBind
	di.MQTTBind
	//di.ProbesBind
	di.WidgetBind
	di.WorkersBind

	client *apcupsd.Client
}

const (
	VariableAPC                         = ""
	VariableDate                        = ""
	VariableHostname                    = ""
	VariableUPSName                     = "ups.model"
	VariableVersion                     = "server.version"
	VariableCable                       = ""
	VariableModel                       = "ups.model"
	VariableUPSMode                     = ""
	VariableStartTime                   = ""
	VariableEndAPC                      = ""
	VariableStatus                      = "ups.status"
	VariableLineVoltage                 = "input.voltage"
	VariableLoadPercent                 = "load"
	VariableBatteryChargePercent        = "battery.charge"
	VariableTimeLeft                    = "battery.runtime"
	VariableMinimumBatteryChargePercent = "battery.charge.low"
	VariableMinimumTimeLeft             = ""
	VariableMaximumTime                 = ""
	VariableMaxLineVoltage              = "input.voltage.maximum"
	VariableMinLineVoltage              = "input.voltage.minimum"
	VariableOutputVoltage               = "output.voltage"
	VariableSense                       = "input.sensitivity"
	VariableDelayWake                   = ""
	VariableDelayShutdown               = "ups.delay.shutdown"
	VariableDelayLowBattery             = "battery.runtime.low"
	VariableLowTransferVoltage          = "input.transfer.low"
	VariableHighTransferVoltage         = "input.transfer.high"
	VariableReturnPercent               = ""
	VariableInternalTemp                = "ups.temperature"
	VariableAlarmDelay                  = ""
	VariableBatteryVoltage              = "battery.voltage"
	VariableLineFrequency               = "input.frequency"
	VariableLastTransfer                = "input.transfer.reason"
	VariableTransfers                   = ""
	VariableXOnBattery                  = ""
	VariableTimeOnBattery               = ""
	VariableCumulativeTimeOnBattery     = ""
	VariableXOffBattery                 = ""
	VariableSelfTest                    = "ups.test.result"
	VariableSelfTestInterval            = "ups.test.interval"
	VariableStatusFlags                 = ""
	VariableDipSwitch                   = ""
	VariableFaultRegister1              = ""
	VariableFaultRegister2              = ""
	VariableFaultRegister3              = ""
	VariableManufacturedDate            = "ups.mfr.date"
	VariableSerialNumber                = "device.serial"
	VariableBatteryDate                 = "battery.date"
	VariableNominalOutputVoltage        = "output.voltage.nominal"
	VariableNominalInputVoltage         = "input.voltage.nominal"
	VariableNominalBatteryVoltage       = "battery.voltage.nominal"
	VariableNominalPower                = "ups.realpower.nominal"
	VariableHumidity                    = "ambient.1.humidity"
	VariableAmbientTemperature          = ""
	VariableExternalBatteries           = "battery.packs.external"
	VariableBadBatteryPacks             = "battery.packs.bad"
	VariableFirmware                    = "ups.firmware"
	VariableAPCModel                    = "ups.model"
)

func (b *Bind) config() *Config {
	return b.Config().Bind().(*Config)
}

func (b *Bind) Run() error {
	cfg := b.config()

	var statusReader, eventsReader apcupsd.Adapter

	if cfg.StatusFile != "" {
		statusReader = file.NewStatusReader(cfg.StatusFile)
	}

	if statusReader == nil && cfg.Address.Host != "" {
		statusReader = nis.NewStatusReader(cfg.Address.Host)
	}

	if statusReader == nil {
		return errors.New("status reader isn't initialized")
	}

	if cfg.EventsFile != "" {
		eventsReader = file.NewStatusReader(cfg.EventsFile)
	}

	if eventsReader == nil && cfg.Address.Host != "" {
		eventsReader = nis.NewEventsReader(cfg.Address.Host)
	}

	if eventsReader == nil {
		return errors.New("events reader isn't initialized")
	}

	b.client = apcupsd.NewClient(statusReader, eventsReader)

	return nil
}
