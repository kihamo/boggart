#!/usr/bin/env bash

influxd backup -database metrics `date +/backups/influx_metrics_%Y%m%d`

# DROP MEASUREMENT boggart_device_heat_meter_pulsar_consumption_cubic_metres_per_hour
# SELECT * INTO boggart_device_heat_meter_pulsar_consumption_cubic_metres_per_hour FROM boggart_pulsar_consumption_cubic_metres_per_hour
# DROP MEASUREMENT boggart_pulsar_consumption_cubic_metres_per_hour