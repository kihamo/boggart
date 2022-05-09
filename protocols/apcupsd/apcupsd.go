package apcupsd

import (
	"context"
	"io"
	"time"
)

type StatusStatus struct {
	AsString string

	IsCal               bool
	IsTrim              bool
	IsBoost             bool
	IsOnline            bool
	IsOnBattery         bool
	IsOverload          bool
	IsLowBattery        bool
	IsReplaceBattery    bool
	IsNoBattery         bool
	IsSlave             bool
	IsSlaveDown         bool
	IsCommunicationLost bool
	IsShuttingDown      bool
}

func (s StatusStatus) String() string {
	return s.AsString
}

type Status struct {
	APC                         string
	Date                        time.Time
	Hostname                    string
	UPSName                     string
	Version                     string
	Cable                       string
	Model                       string
	UPSMode                     string
	StartTime                   time.Time
	EndAPC                      time.Time
	Status                      StatusStatus
	LineVoltage                 *float64
	LoadPercent                 *float64
	BatteryChargePercent        *float64
	TimeLeft                    *time.Duration
	MinimumBatteryChargePercent *float64
	MinimumTimeLeft             *time.Duration
	MaximumTime                 *time.Duration
	MaxLineVoltage              *float64
	MinLineVoltage              *float64
	OutputVoltage               *float64
	Sense                       *string
	DelayWake                   *time.Duration
	DelayShutdown               *time.Duration
	DelayLowBattery             *time.Duration
	LowTransferVoltage          *float64
	HighTransferVoltage         *float64
	ReturnPercent               *float64
	InternalTemp                *float64
	AlarmDelay                  *time.Duration
	BatteryVoltage              *float64
	LineFrequency               *float64
	LastTransfer                *string
	Transfers                   *uint64
	XOnBattery                  *time.Time
	TimeOnBattery               *time.Duration
	CumulativeTimeOnBattery     *time.Duration
	XOffBattery                 *time.Time
	SelfTest                    *string
	SelfTestInterval            *time.Duration
	StatusFlags                 *string
	DipSwitch                   *string
	FaultRegister1              *string
	FaultRegister2              *string
	FaultRegister3              *string
	ManufacturedDate            *time.Time
	SerialNumber                *string
	BatteryDate                 *time.Time
	NominalOutputVoltage        *float64
	NominalInputVoltage         *float64
	NominalBatteryVoltage       *float64
	NominalPower                *float64
	Humidity                    *float64
	AmbientTemperature          *float64
	ExternalBatteries           *uint64
	BadBatteryPacks             *uint64
	Firmware                    *string
	APCModel                    *string
}

type Event struct {
	Date    time.Time
	Message string
}

type Adapter interface {
	Reader(context.Context) (io.Reader, error)
}
