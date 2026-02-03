package mqtt

import (
	"context"
	"net/http"
	"net/url"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/snitch"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func (b *Bind) Tasks() []tasks.Task {
	cfg := b.config()

	if cfg.IPAddressSensorID == "" {
		return nil
	}

	return []tasks.Task{
		tasks.NewTask().
			WithName("import-metrics").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskImportMetricsHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.ImportMetricsInterval)),
	}
}

func (b *Bind) taskImportMetricsHandler(ctx context.Context) error {
	mac := b.Meta().MACAsString()
	if mac == "" {
		return nil
	}

	ip := b.IP()
	if ip == nil {
		return nil
	}

	u := &url.URL{
		Scheme: "http",
		Host:   ip.String(),
		Path:   "/metrics",
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	parser := expfmt.TextParser{}
	m, err := parser.TextToMetricFamilies(response.Body)
	if err != nil {
		return err
	}

	for key, outerMetric := range m {
		t := outerMetric.GetType()
		if t != io_prometheus_client.MetricType_GAUGE {
			b.Logger().Warn("Unknown metric type", "key", key, "type", t.String())
			continue
		}

		var innerMetric snitch.Gauge

		switch key {
		case "esphome_sensor_value":
			innerMetric = metricSensorValue
		case "esphome_sensor_failed":
			innerMetric = metricSensorFailed
		case "esphome_binary_sensor_value":
			innerMetric = metricBinarySensorValue
		case "esphome_binary_sensor_failed":
			innerMetric = metricBinarySensorFailed
		case "esphome_fan_value":
			innerMetric = metricFanValue
		case "esphome_fan_failed":
			innerMetric = metricFanFailed
		case "esphome_fan_speed":
			innerMetric = metricFanSpeed
		case "esphome_fan_oscillation":
			innerMetric = metricFanOscillation
		case "esphome_light_state":
			innerMetric = metricLightState
		case "esphome_light_color":
			innerMetric = metricLightColor
		case "esphome_light_effect_active":
			innerMetric = metricLightEffectActive
		case "esphome_cover_value":
			innerMetric = metricCoverValue
		case "esphome_cover_failed":
			innerMetric = metricCoverFailed
		case "esphome_switch_value":
			innerMetric = metricSwitchValue
		case "esphome_switch_failed":
			innerMetric = metricSwitchFailed
		case "esphome_lock_value":
			innerMetric = metricLockValue
		case "esphome_lock_failed":
			innerMetric = metricLockFailed
		case "esphome_text_sensor_value":
			innerMetric = metricTextSensorValue
		case "esphome_text_sensor_failed":
			innerMetric = metricTextSensorFailed
		case "esphome_text_value":
			innerMetric = metricTextValue
		case "esphome_text_failed":
			innerMetric = metricTextFailed
		case "esphome_event_value":
			innerMetric = metricEventValue
		case "esphome_event_failed":
			innerMetric = metricEventFailed
		case "esphome_number_value":
			innerMetric = metricNumberValue
		case "esphome_number_failed":
			innerMetric = metricNumberFailed
		case "esphome_select_value":
			innerMetric = metriceSelectValue
		case "esphome_select_failed":
			innerMetric = metricSelectFailed
		case "esphome_media_player_state_value":
			innerMetric = metricMediaPlayerStateValue
		case "esphome_media_player_volume":
			innerMetric = metricMediaPlayerVolume
		case "esphome_media_player_is_muted":
			innerMetric = metricMediaPlayerIsMuted
		case "esphome_media_player_failed":
			innerMetric = metricMediaPlayerFailed
		case "esphome_update_entity_state":
			innerMetric = metricUpdateEntityState
		case "esphome_update_entity_info":
			innerMetric = metricUpdateEntityInfo
		case "esphome_update_entity_failed":
			innerMetric = metricUpdateEntityFailed
		case "esphome_valve_operation":
			innerMetric = metricValveOperation
		case "esphome_valve_failed":
			innerMetric = metricValveFailed
		case "esphome_valve_position":
			innerMetric = metricValvePosition
		case "esphome_climate_setting":
			innerMetric = metricClimateSetting
		case "esphome_climate_value":
			innerMetric = metricClimateValue
		case "esphome_climate_failed":
			innerMetric = metricClimateFailed
		}

		if innerMetric == nil {
			b.Logger().Warn("Unknown metric key", "key", key, "type", t.String())
			continue
		}

		innerMetric = innerMetric.With("mac", mac)

		for _, subMetric := range outerMetric.GetMetric() {
			gauge := subMetric.GetGauge()
			if gauge == nil {
				b.Logger().Warn("Gauge metric is empty", "key", key)
				continue
			}

			outerLabels := subMetric.GetLabel()
			innerLabels := make([]string, 0, len(outerLabels)*2)

			for _, label := range outerLabels {
				innerLabels = append(innerLabels, label.GetName(), label.GetValue())
			}

			innerMetric.With(innerLabels...).Set(gauge.GetValue())
		}
	}

	return nil
}
