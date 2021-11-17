import esphome.codegen as cg
import esphome.config_validation as cv
from esphome.components import sensor, uart
from esphome.const import (
    CONF_ID,
    CONF_ADDRESS,
    CONF_POWER,
    CONF_VOLTAGE,
    DEVICE_CLASS_CURRENT,
    DEVICE_CLASS_ENERGY,
    DEVICE_CLASS_POWER,
    DEVICE_CLASS_VOLTAGE,
    ICON_FLASH,
    STATE_CLASS_MEASUREMENT,
    STATE_CLASS_TOTAL_INCREASING,
    UNIT_VOLT,
    UNIT_AMPERE,
    UNIT_WATT,
    UNIT_WATT_HOURS,
)

DEPENDENCIES = ["uart"]

CONF_AMPERAGE = "amperage"
CONF_TARIFF1 = "tariff1"
CONF_TARIFF2 = "tariff2"
CONF_TARIFF3 = "tariff3"
CONF_TARIFF4 = "tariff4"
CONF_TARIFFS_TOTAL = "tariffs_total"

mercury1_ns = cg.esphome_ns.namespace("mercury1")
Mercury1 = mercury1_ns.class_("Mercury1", cg.PollingComponent, uart.UARTDevice)

CONFIG_SCHEMA = (
    cv.Schema(
        {
            cv.GenerateID(): cv.declare_id(Mercury1),
            cv.Required(CONF_ADDRESS): cv.hex_uint32_t,
            cv.Optional(CONF_VOLTAGE): sensor.sensor_schema(
                unit_of_measurement=UNIT_VOLT,
                accuracy_decimals=1,
                device_class=DEVICE_CLASS_VOLTAGE,
                state_class=STATE_CLASS_MEASUREMENT,
                icon=ICON_FLASH,
            ),
            cv.Optional(CONF_AMPERAGE): sensor.sensor_schema(
                unit_of_measurement=UNIT_AMPERE,
                accuracy_decimals=2,
                device_class=DEVICE_CLASS_CURRENT,
                state_class=STATE_CLASS_MEASUREMENT,
                icon=ICON_FLASH,
            ),
            cv.Optional(CONF_POWER): sensor.sensor_schema(
                unit_of_measurement=UNIT_WATT,
                accuracy_decimals=0,
                device_class=DEVICE_CLASS_POWER,
                state_class=STATE_CLASS_MEASUREMENT,
                icon=ICON_FLASH,
            ),
            cv.Optional(CONF_TARIFF1): sensor.sensor_schema(
                unit_of_measurement=UNIT_WATT_HOURS,
                accuracy_decimals=0,
                device_class=DEVICE_CLASS_ENERGY,
                state_class=STATE_CLASS_TOTAL_INCREASING,
                icon=ICON_FLASH,
            ),
            cv.Optional(CONF_TARIFF2): sensor.sensor_schema(
                unit_of_measurement=UNIT_WATT_HOURS,
                accuracy_decimals=0,
                device_class=DEVICE_CLASS_ENERGY,
                state_class=STATE_CLASS_TOTAL_INCREASING,
                icon=ICON_FLASH,
            ),
            cv.Optional(CONF_TARIFF3): sensor.sensor_schema(
                unit_of_measurement=UNIT_WATT_HOURS,
                accuracy_decimals=0,
                device_class=DEVICE_CLASS_ENERGY,
                state_class=STATE_CLASS_TOTAL_INCREASING,
                icon=ICON_FLASH,
            ),
            cv.Optional(CONF_TARIFF4): sensor.sensor_schema(
                unit_of_measurement=UNIT_WATT_HOURS,
                accuracy_decimals=0,
                device_class=DEVICE_CLASS_ENERGY,
                state_class=STATE_CLASS_TOTAL_INCREASING,
                icon=ICON_FLASH,
            ),
            cv.Optional(CONF_TARIFFS_TOTAL): sensor.sensor_schema(
                unit_of_measurement=UNIT_WATT_HOURS,
                accuracy_decimals=0,
                device_class=DEVICE_CLASS_ENERGY,
                state_class=STATE_CLASS_TOTAL_INCREASING,
                icon=ICON_FLASH,
            ),
        }
    )
        .extend(cv.polling_component_schema("60s"))
        .extend(uart.UART_DEVICE_SCHEMA)
)

async def to_code(config):
    var = cg.new_Pvariable(config[CONF_ID])
    cg.add(var.set_address(config[CONF_ADDRESS]))

    await cg.register_component(var, config)
    await uart.register_uart_device(var, config)

    if CONF_VOLTAGE in config:
        conf = config[CONF_VOLTAGE]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_voltage_sensor(sens))
    if CONF_AMPERAGE in config:
        conf = config[CONF_AMPERAGE]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_amperage_sensor(sens))
    if CONF_POWER in config:
        conf = config[CONF_POWER]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_power_sensor(sens))
    if CONF_TARIFF1 in config:
        conf = config[CONF_TARIFF1]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_tariff1_sensor(sens))
    if CONF_TARIFF2 in config:
        conf = config[CONF_TARIFF2]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_tariff2_sensor(sens))
    if CONF_TARIFF3 in config:
        conf = config[CONF_TARIFF3]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_tariff3_sensor(sens))
    if CONF_TARIFF4 in config:
        conf = config[CONF_TARIFF4]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_tariff4_sensor(sens))
    if CONF_TARIFFS_TOTAL in config:
        conf = config[CONF_TARIFFS_TOTAL]
        sens = await sensor.new_sensor(conf)
        cg.add(var.set_tariffs_total_sensor(sens))
