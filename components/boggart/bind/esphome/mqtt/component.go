package mqtt

import (
	"github.com/kihamo/boggart/components/mqtt"
)

type ComponentType string

const (
	ComponentTypeUnknown      ComponentType = "unknown"
	ComponentTypeBinarySensor ComponentType = "binary_sensor"
	ComponentTypeCover        ComponentType = "cover"
	ComponentTypeFan          ComponentType = "fan"
	ComponentTypeLight        ComponentType = "light"
	ComponentTypeSensor       ComponentType = "sensor"
	ComponentTypeSwitch       ComponentType = "switch"
	ComponentTypeTextSensor   ComponentType = "text_sensor"
	ComponentTypeCamera       ComponentType = "camera"
	ComponentTypeClimate      ComponentType = "climate"

	// https://www.home-assistant.io/integrations/binary_sensor
	// https://www.home-assistant.io/integrations/sensor#device-class
	DeviceClassBatteryCharging          = "battery_charging"
	DeviceClassCold                     = "cold"
	DeviceClassConnectivity             = "connectivity"
	DeviceClassDoor                     = "door"
	DeviceClassGarageDoor               = "garage_door"
	DeviceClassHeat                     = "heat"
	DeviceClassLight                    = "light"
	DeviceClassLock                     = "lock"
	DeviceClassMoisture                 = "moisture"
	DeviceClassMotion                   = "motion"
	DeviceClassMoving                   = "moving"
	DeviceClassOccupancy                = "occupancy"
	DeviceClassOpening                  = "opening"
	DeviceClassPlug                     = "plug"
	DeviceClassPresence                 = "presence"
	DeviceClassProblem                  = "problem"
	DeviceClassRunning                  = "running"
	DeviceClassSafety                   = "safety"
	DeviceClassSmoke                    = "smoke"
	DeviceClassSound                    = "sound"
	DeviceClassTamper                   = "tamper"
	DeviceClassUpdate                   = "update"
	DeviceClassVibration                = "vibration"
	DeviceClassWindow                   = "window"
	DeviceClassEmpty                    = ""
	DeviceClassBattery                  = "battery"                    // Percentage of battery that is left.
	DeviceClassQas                      = "gas"                        // Gasvolume in m³ or ft³.
	DeviceClassPower                    = "power"                      // Power in W or kW.
	DeviceClassAQI                      = "aqi"                        // Air Quality Index
	DeviceClassCarbonDioxide            = "carbon_dioxide"             // Carbon Dioxide in CO2 (Smoke)
	DeviceClassCarbonMonoxide           = "carbon_monoxide"            // Carbon Monoxide in CO (Gas CNG/LPG)
	DeviceClassCurrent                  = "current"                    // Current in A.
	DeviceClassEnergy                   = "energy"                     // Energy in Wh or kWh.
	DeviceClassHumidity                 = "humidity"                   // Percentage of humidity in the air.
	DeviceClassIlluminance              = "illuminance"                // The current light level in lx or lm.
	DeviceClassMonetary                 = "monetary"                   // The monetary value.
	DeviceClassNitrogenDioxide          = "nitrogen_dioxide"           // Concentration of Nitrogen Dioxide in µg/m³
	DeviceClassNitrogenMonoxide         = "nitrogen_monoxide"          // Concentration of Nitrogen Monoxide in µg/m³
	DeviceClassNitrousOxide             = "nitrous_oxide"              // Concentration of Nitrous Oxide in µg/m³
	DeviceClassOzone                    = "ozone"                      // Concentration of Ozone in µg/m³
	DeviceClassPm1                      = "pm1"                        // Concentration of particulate matter less than 1 micrometer in µg/m³
	DeviceClassPm10                     = "pm10"                       // Concentration of particulate matter less than 10 micrometers in µg/m³
	DeviceClassPm25                     = "pm25"                       // Concentration of particulate matter less than 2.5 micrometers in µg/m³
	DeviceClassPowerFactor              = "power_factor"               // Power factor in %.
	DeviceClassPressure                 = "pressure"                   // Pressure in hPa or mbar.
	DeviceClassSignalStrength           = "signal_strength"            // Signal strength in dB or dBm.
	DeviceClassSulphurDioxide           = "sulphur_dioxide"            // Concentration of sulphur dioxide in µg/m³
	DeviceClassTemperature              = "temperature"                // Temperature in °C or °F.
	DeviceClassDate                     = "date"                       // Date string (ISO 8601).
	DeviceClassTimestamp                = "timestamp"                  // Datetime object or timestamp string (ISO 8601).
	DeviceClassVolatileOrganicCompounds = "volatile_organic_compounds" // Concentration of volatile organic compounds in µg/m³.
	DeviceClassVoltage                  = "voltage"                    // Voltage in V.
)

func (t ComponentType) String() string {
	return string(t)
}

type Component interface {
	ID() string
	Type() ComponentType
	UniqueID() string
	Name() string
	Icon() string
	State() interface{}
	StateFormat() string
	SetState(mqtt.Message) error
	ConfigMessage() mqtt.Message
	StateTopic() mqtt.Topic
	CommandTopic() mqtt.Topic
	AvailabilityTopic() mqtt.Topic
	DeviceInfo() DeviceInfo
	CommandToPayload(cmd interface{}) interface{}
	Subscribe(subscribers ...mqtt.Subscriber)
	Subscribers() []mqtt.Subscriber
}
