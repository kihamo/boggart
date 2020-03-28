package apcupsd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	timeFormatLong  = "2006-01-02 15:04:05 -0700"
	timeFormatShort = "2006-01-02"
)

// http://www.apcupsd.org/manual/manual.html#status-report-fields

const (
	StatusFieldAPC                         = "APC"       // Header record indicating the STATUS format revision level, the number of records that follow the APC statement, and the number of bytes that follow the record.
	StatusFieldDate                        = "DATE"      // The date and time that the information was last obtained from the UPS.
	StatusFieldHostname                    = "HOSTNAME"  // The name of the machine that collected the UPS data.
	StatusFieldUPSName                     = "UPSNAME"   // The name of the UPS as stored in the EEPROM or in the UPSNAME directive in the configuration file.
	StatusFieldVersion                     = "VERSION"   // The apcupsd release number, build date, and platform.
	StatusFieldCable                       = "CABLE"     // The cable as specified in the configuration file (UPSCABLE).
	StatusFieldModel                       = "MODEL"     // The UPS model as derived from information from the UPS.
	StatusFieldUPSMode                     = "UPSMODE"   // The mode in which apcupsd is operating as specified in the configuration file (UPSMODE)
	StatusFieldStartTime                   = "STARTTIME" // The time/date that apcupsd was started.
	StatusFieldStatus                      = "STATUS"    // The current status of the UPS (ONLINE, ONBATT, etc.)
	StatusFieldLineVoltage                 = "LINEV"     // The current line voltage as returned by the UPS.
	StatusFieldLoadPercent                 = "LOADPCT"   // The percentage of load capacity as estimated by the UPS.
	StatusFieldBatteryChargePercent        = "BCHARGE"   // The percentage charge on the batteries.
	StatusFieldTimeLeft                    = "TIMELEFT"  // The remaining runtime left on batteries as estimated by the UPS.
	StatusFieldMinimumBatteryChargePercent = "MBATTCHG"  // If the battery charge percentage (BCHARGE) drops below this value, apcupsd will shutdown your system. Value is set in the configuration file (BATTERYLEVEL)
	StatusFieldMinimumTimeLeft             = "MINTIMEL"  // apcupsd will shutdown your system if the time on batteries exceeds this value. A value of zero disables the feature. Value is set in the configuration file (TIMEOUT)
	StatusFieldMaximumTime                 = "MAXTIME"   // apcupsd will shutdown your system if the time on batteries exceeds this value. A value of zero disables the feature. Value is set in the configuration file (TIMEOUT)
	StatusFieldMaxLineVoltage              = "MAXLINEV"  // The maximum line voltage since the UPS was started, as reported by the UPS
	StatusFieldMinLineVoltage              = "MINLINEV"  // The minimum line voltage since the UPS was started, as returned by the UPS
	StatusFieldOutputVoltage               = "OUTPUTV"   // The voltage the UPS is supplying to your equipment
	StatusFieldSense                       = "SENSE"     // The sensitivity level of the UPS to line voltage fluctuations.
	StatusFieldDelayWake                   = "DWAKE"     // The amount of time the UPS will wait before restoring power to your equipment after a power off condition when the power is restored.
	StatusFieldDelayShutdown               = "DSHUTD"    // The grace delay that the UPS gives after receiving a power down command from apcupsd before it powers off your equipment.
	StatusFieldDelayLowBattery             = "DLOWBATT"  // The remaining runtime below which the UPS sends the low battery signal. At this point apcupsd will force an immediate emergency shutdown.
	StatusFieldLowTransferVoltage          = "LOTRANS"   // The line voltage below which the UPS will switch to batteries
	StatusFieldHighTransferVoltage         = "HITRANS"   // The line voltage above which the UPS will switch to batteries.
	StatusFieldReturnPercent               = "RETPCT"    // The percentage charge that the batteries must have after a power off condition before the UPS will restore power to your equipment.
	StatusFieldInternalTemp                = "ITEMP"     // Internal UPS temperature as supplied by the UPS.
	StatusFieldAlarmDelay                  = "ALARMDEL"  // The delay period for the UPS alarm.
	StatusFieldBatteryVoltage              = "BATTV"     // Battery voltage as supplied by the UPS.
	StatusFieldLineFrequency               = "LINEFREQ"  // Line frequency in hertz as given by the UPS.
	StatusFieldLastTransfer                = "LASTXFER"  // The reason for the last transfer to batteries.
	StatusFieldTransfers                   = "NUMXFERS"  // The number of transfers to batteries since apcupsd startup.
	StatusFieldXOnBattery                  = "XONBATT"   // Time and date of last transfer to batteries, or N/A.
	StatusFieldTimeOnBattery               = "TONBATT"   // Time in seconds currently on batteries, or 0.
	StatusFieldCumulativeTimeOnBattery     = "CUMONBATT" // Total (cumulative) time on batteries in seconds since apcupsd startup.
	StatusFieldXOffBattery                 = "XOFFBATT"  // Time and date of last transfer from batteries, or N/A.
	StatusFieldSelfTest                    = "SELFTEST"  // Self test interval in hours (336 = 2 weeks, 168 = 1 week, ON = at power on, OFF = never).
	StatusFieldSelfTestInterval            = "STESTI"    // The interval in hours between automatic self tests.
	StatusFieldStatusFlags                 = "STATFLAG"  // Status flag. English version is given by STATUS.
	StatusFieldDipSwitch                   = "DIPSW"     // The current dip switch settings on UPSes that have them
	StatusFieldFaultRegister1              = "REG1"      // The value from the UPS fault register 1.
	StatusFieldFaultRegister2              = "REG2"      // The value from the UPS fault register 2.
	StatusFieldFaultRegister3              = "REG3"      // The value from the UPS fault register 3.
	StatusFieldManufacturedDate            = "MANDATE"   // The date the UPS was manufactured.
	StatusFieldSerialNumber                = "SERIALNO"  // The UPS serial number.
	StatusFieldBatteryDate                 = "BATTDATE"  // The date that batteries were last replaced.
	StatusFieldNominalOutputVoltage        = "NOMOUTV"   // The output voltage that the UPS will attempt to supply when on battery power.
	StatusFieldNominalInputVoltage         = "NOMINV"    // The input voltage that the UPS is configured to expect.
	StatusFieldNominalBatteryVoltage       = "NOMBATTV"  // The nominal battery voltage.
	StatusFieldNominalPower                = "NOMPOWER"  // The maximum power in Watts that the UPS is designed to supply.
	StatusFieldHumidity                    = "HUMIDITY"  // The humidity as measured by the UPS.
	StatusFieldAmbientTemperature          = "AMBTEMP"   // The ambient temperature as measured by the UPS.
	StatusFieldExternalBatteries           = "EXTBATTS"  // The number of external batteries as defined by the user. A correct number here helps the UPS compute the remaining runtime more accurately.
	StatusFieldBadBatteryPacks             = "BADBATTS"  // The number of bad battery packs.
	StatusFieldFirmware                    = "FIRMWARE"  // The firmware revision number as reported by the UPS.
	StatusFieldAPCModel                    = "APCMODEL"  // The old APC model identification code.
	StatusFieldEndAPC                      = "END APC"   // The time and date that the STATUS record was written.

	StatusStatusCal               = "CAL"
	StatusStatusTrim              = "TRIM"
	StatusStatusBoost             = "BOOST"
	StatusStatusOnline            = "ONLINE"
	StatusStatusOnBattery         = "ONBATT"
	StatusStatusOverload          = "OVERLOAD"
	StatusStatusLowBattery        = "LOWBATT"
	StatusStatusReplaceBattery    = "REPLACEBATT"
	StatusStatusNoBattery         = "NOBATT"
	StatusStatusSlave             = "SLAVE"
	StatusStatusSlaveDown         = "SLAVEDOWN"
	StatusStatusCommunicationLost = "COMMLOST"
	StatusStatusShuttingDown      = "SHUTTING"
)

type Client struct {
	status Adapter
	events Adapter
}

func NewClient(status Adapter, events Adapter) *Client {
	return &Client{
		status: status,
		events: events,
	}
}

func (c *Client) Status(ctx context.Context) (Status, error) {
	status := Status{}

	reader, err := c.status.Reader(ctx)
	if err != nil {
		return status, err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		sp := strings.SplitN(scanner.Text(), ":", 2)
		if len(sp) != 2 {
			return status, fmt.Errorf("invalid key/value pair for string '%s'", scanner.Text())
		}

		var err error

		key := strings.TrimSpace(sp[0])
		value := strings.TrimSpace(sp[1])

		switch key {
		case StatusFieldAPC:
			status.APC = ParseString(value)
		case StatusFieldDate:
			status.Date, err = ParseTime(timeFormatLong, value)
		case StatusFieldHostname:
			status.Hostname = ParseString(value)
		case StatusFieldUPSName:
			status.UPSName = ParseString(value)
		case StatusFieldVersion:
			status.Version = ParseString(value)
		case StatusFieldCable:
			status.Cable = ParseString(value)
		case StatusFieldModel:
			status.Model = ParseString(value)
		case StatusFieldUPSMode:
			status.UPSMode = ParseString(value)
		case StatusFieldStartTime:
			status.StartTime, err = ParseTime(timeFormatLong, value)
		case StatusFieldStatus:
			s := ParseString(value)
			if s != nil && *s != "" {
				status.Status = &StatusStatus{
					AsString: *s,
				}

				for _, name := range strings.Fields(*s) {
					switch name {
					case StatusStatusCal:
						status.Status.IsCal = true
					case StatusStatusTrim:
						status.Status.IsTrim = true
					case StatusStatusBoost:
						status.Status.IsBoost = true
					case StatusStatusOnline:
						status.Status.IsOnline = true
					case StatusStatusOnBattery:
						status.Status.IsOnBattery = true
					case StatusStatusOverload:
						status.Status.IsOverload = true
					case StatusStatusLowBattery:
						status.Status.IsLowBattery = true
					case StatusStatusReplaceBattery:
						status.Status.IsReplaceBattery = true
					case StatusStatusNoBattery:
						status.Status.IsNoBattery = true
					case StatusStatusSlave:
						status.Status.IsSlave = true
					case StatusStatusSlaveDown:
						status.Status.IsSlaveDown = true
					case StatusStatusCommunicationLost:
						status.Status.IsCommunicationLost = true
					case StatusStatusShuttingDown:
						status.Status.IsShuttingDown = true
					}
				}
			}

		case StatusFieldLineVoltage:
			status.LineVoltage, err = ParseFloat(value)
		case StatusFieldLoadPercent:
			status.LoadPercent, err = ParseFloat(value)
		case StatusFieldBatteryChargePercent:
			status.BatteryChargePercent, err = ParseFloat(value)
		case StatusFieldTimeLeft:
			status.TimeLeft, err = ParseDuration(value)
		case StatusFieldMinimumBatteryChargePercent:
			status.MinimumBatteryChargePercent, err = ParseFloat(value)
		case StatusFieldMinimumTimeLeft:
			status.MinimumTimeLeft, err = ParseDuration(value)
		case StatusFieldMaximumTime:
			status.MaximumTime, err = ParseDuration(value)
		case StatusFieldMaxLineVoltage:
			status.MaxLineVoltage, err = ParseFloat(value)
		case StatusFieldMinLineVoltage:
			status.MinLineVoltage, err = ParseFloat(value)
		case StatusFieldOutputVoltage:
			status.OutputVoltage, err = ParseFloat(value)
		case StatusFieldSense:
			status.Sense = ParseString(value)
		case StatusFieldDelayWake:
			status.DelayWake, err = ParseDuration(value)
		case StatusFieldDelayShutdown:
			status.DelayShutdown, err = ParseDuration(value)
		case StatusFieldDelayLowBattery:
			status.DelayLowBattery, err = ParseDuration(value)
		case StatusFieldLowTransferVoltage:
			status.LowTransferVoltage, err = ParseFloat(value)
		case StatusFieldHighTransferVoltage:
			status.HighTransferVoltage, err = ParseFloat(value)
		case StatusFieldReturnPercent:
			status.ReturnPercent, err = ParseFloat(value)
		case StatusFieldInternalTemp:
			status.InternalTemp, err = ParseFloat(value)
		case StatusFieldAlarmDelay:
			if value == "No alarm" {
				break
			}

			status.AlarmDelay, err = ParseDuration(value)
		case StatusFieldBatteryVoltage:
			status.BatteryVoltage, err = ParseFloat(value)
		case StatusFieldLineFrequency:
			status.LineFrequency, err = ParseFloat(value)
		case StatusFieldLastTransfer:
			status.LastTransfer = ParseString(value)
		case StatusFieldTransfers:
			status.Transfers, err = ParseUint(value)
		case StatusFieldXOnBattery:
			status.XOnBattery, err = ParseTime(timeFormatLong, value)
		case StatusFieldTimeOnBattery:
			status.TimeOnBattery, err = ParseDuration(value)
		case StatusFieldCumulativeTimeOnBattery:
			status.CumulativeTimeOnBattery, err = ParseDuration(value)
		case StatusFieldXOffBattery:
			status.XOffBattery, err = ParseTime(timeFormatLong, value)
		case StatusFieldSelfTest:
			status.SelfTest = ParseString(value)
		case StatusFieldSelfTestInterval:
			if value != "" {
				value += " Hours"
			}

			status.SelfTest = ParseString(value)
		case StatusFieldStatusFlags:
			status.StatusFlags = ParseString(value)
		case StatusFieldDipSwitch:
			status.DipSwitch = ParseString(value)
		case StatusFieldFaultRegister1:
			status.FaultRegister1 = ParseString(value)
		case StatusFieldFaultRegister2:
			status.FaultRegister2 = ParseString(value)
		case StatusFieldFaultRegister3:
			status.FaultRegister3 = ParseString(value)
		case StatusFieldManufacturedDate:
			status.ManufacturedDate, err = ParseTime(timeFormatShort, value)
		case StatusFieldSerialNumber:
			status.SerialNumber = ParseString(value)
		case StatusFieldBatteryDate:
			status.BatteryDate, err = ParseTime(timeFormatShort, value)
		case StatusFieldNominalOutputVoltage:
			status.NominalOutputVoltage, err = ParseFloat(value)
		case StatusFieldNominalInputVoltage:
			status.NominalInputVoltage, err = ParseFloat(value)
		case StatusFieldNominalBatteryVoltage:
			status.NominalBatteryVoltage, err = ParseFloat(value)
		case StatusFieldNominalPower:
			status.NominalPower, err = ParseFloat(value)
		case StatusFieldHumidity:
			status.NominalPower, err = ParseFloat(value)
		case StatusFieldAmbientTemperature:
			status.AmbientTemperature, err = ParseFloat(value)
		case StatusFieldExternalBatteries:
			status.ExternalBatteries, err = ParseUint(value)
		case StatusFieldBadBatteryPacks:
			status.BadBatteryPacks, err = ParseUint(value)
		case StatusFieldFirmware:
			status.Firmware = ParseString(value)
		case StatusFieldAPCModel:
			status.APCModel = ParseString(value)
		case StatusFieldEndAPC:
			status.EndAPC, err = ParseTime(timeFormatLong, value)
		}

		if err != nil {
			return status, fmt.Errorf("Failed parse %s key with error %s", key, err.Error())
		}
	}

	if err := scanner.Err(); err != nil {
		return status, err
	}

	return status, nil
}

func (c *Client) Events(ctx context.Context) ([]Event, error) {
	reader, err := c.events.Reader(ctx)
	if err != nil {
		return nil, err
	}

	events := make([]Event, 0)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()

		if len(text) < len(timeFormatLong)+2 {
			return nil, errors.New("Event string is short")
		}

		t, err := ParseTime(timeFormatLong, text[:len(timeFormatLong)])
		if err != nil {
			return nil, err
		}

		events = append([]Event{{
			Date:    *t,
			Message: text[len(timeFormatLong)+2:],
		}}, events...)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func ParseString(value string) *string {
	if value == "N/A" {
		return nil
	}

	return &value
}

func ParseUint(value string) (*uint64, error) {
	if value == "N/A" || value == "" {
		return nil, nil
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 0 || parts[0] == "" {
		return nil, nil
	}

	u, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return nil, nil
	}

	return &u, err
}

func ParseFloat(value string) (*float64, error) {
	if value == "N/A" || value == "" {
		return nil, nil
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) == 0 || parts[0] == "" {
		return nil, nil
	}

	f, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return nil, nil
	}

	return &f, err
}

func ParseTime(layout, value string) (*time.Time, error) {
	if value == "N/A" || value == "" {
		return nil, nil
	}

	t, err := time.Parse(layout, value)
	if err != nil {
		return nil, err
	}

	return &t, err
}

func ParseDuration(value string) (*time.Duration, error) {
	if value == "N/A" || value == "" {
		return nil, nil
	}

	parts := strings.SplitN(value, " ", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid time duration")
	}

	unit := strings.ToLower(string([]rune(parts[1])[0]))

	d, err := time.ParseDuration(fmt.Sprintf("%s%s", parts[0], unit))
	if err != nil {
		return nil, err
	}

	return &d, nil
}
