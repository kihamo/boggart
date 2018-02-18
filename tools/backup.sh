#!/usr/bin/env bash

influxd backup -database metrics `date +/backups/influx_metrics_%Y%m%d`