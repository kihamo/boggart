package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	influxdb "github.com/influxdata/influxdb/client/v2"
)

const (
	DATABASE               = "metrics"
	MEASUREMENT_OLD        = "boggart_pulsar_hot_water_capacity_cubic_metres"
	MEASUREMENT_NEW        = "boggart_device_water_meter_pulsar_pulsed_volume_cubic_metres"
	REMOVE_MEASUREMENT_NEW = false
	COUNT_FIELD            = "value"
	BULK                   = 5000
)

var (
	fieldsToTags = map[string]bool{
		"app_build":     true,
		"app_name":      true,
		"app_version":   true,
		"hostname":      true,
		"serial_number": true,
	}
	fieldsExtends = map[string]interface{}{
		"serial_number": "16_85865",
	}
)

func main() {
	client, err := influxdb.NewHTTPClient(influxdb.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "",
		Password: "",
	})

	if err != nil {
		log.Fatalf("Init client failed with error %s", err.Error())
	}

	response, err := client.Query(influxdb.Query{
		Database: DATABASE,
		Command:  fmt.Sprintf("SELECT count(%s) FROM %s", COUNT_FIELD, MEASUREMENT_OLD),
	})

	if err != nil {
		log.Fatalf("Invalid query execute: %s", err)
	} else {
		log.Printf("Total count in OLD measurement %s: %s\n", MEASUREMENT_OLD, response.Results[0].Series[0].Values[0][1])
	}

	if REMOVE_MEASUREMENT_NEW {
		_, err = client.Query(influxdb.Query{
			Database: DATABASE,
			Command:  fmt.Sprintf("DROP MEASUREMENT %s", MEASUREMENT_NEW),
		})
		if err != nil {
			log.Fatalf("Failed remove new measurement: %s", err)
		} else {
			log.Printf("New measurement %s removed\n", MEASUREMENT_NEW)
		}
	}

	for offset := 0; ; offset += BULK {
		response, err := client.Query(influxdb.Query{
			Database: DATABASE,
			Command:  fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", MEASUREMENT_OLD, BULK, offset),
		})
		if err != nil {
			log.Fatalf("Invalid query execute: %s", err)
		}

		if len(response.Results[0].Series) == 0 {
			log.Printf("END")
			break
		}

		// init batch
		batchPoint, err := influxdb.NewBatchPoints(influxdb.BatchPointsConfig{
			Database:  DATABASE,
			Precision: "s",
		})
		if err != nil {
			log.Fatal(err)
		}

		// data
		for _, row := range response.Results[0].Series[0].Values {
			t, err := time.Parse(time.RFC3339, row[0].(string))
			if err != nil {
				log.Fatal(err)
			}

			tags := map[string]string{}
			fields := map[string]interface{}{}

			for key, value := range fieldsExtends {
				if _, ok := fieldsToTags[key]; ok {
					tags[key] = value.(string)
				} else {
					fields[key] = value
				}
			}

			for i, column := range response.Results[0].Series[0].Columns {
				// skeep time
				if i == 0 {
					continue
				}

				if _, ok := fieldsToTags[column]; ok {
					tags[column] = row[i].(string)
				} else {
					if j, ok := row[i].(json.Number); ok {
						fields[column], _ = j.Float64()
					} else {
						log.Printf("Skip columnt %s with value %t\n", column, row[i])
					}
				}
			}

			point, err := influxdb.NewPoint(MEASUREMENT_NEW, tags, fields, t)
			if err != nil {
				log.Fatal(err)
			}

			batchPoint.AddPoint(point)
			// log.Println("Add point")
		}

		if err := client.Write(batchPoint); err != nil {
			log.Fatal(err)
		}

		log.Printf("LIMIT %d OFFSET %d", BULK, offset)
	}

	response, err = client.Query(influxdb.Query{
		Database: DATABASE,
		Command:  fmt.Sprintf("SELECT count(%s) FROM %s", COUNT_FIELD, MEASUREMENT_NEW),
	})
	if err != nil {
		log.Fatalf("Invalid query execute: %s", err)
	} else {
		log.Printf("Total count in NEW measurement %s: %s\n", MEASUREMENT_NEW, response.Results[0].Series[0].Values[0][1])
	}
}
