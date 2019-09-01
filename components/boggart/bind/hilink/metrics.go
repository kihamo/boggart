package hilink

import (
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/snitch"
)

var (
	metricBalance              = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_balance_rubles", "HiLink balance in rubles")
	metricLimitInternetTraffic = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_limit_internet_traffic_bytes", "HiLink limit internet traffic in bytes")
	metricSignalRSSI           = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_signal_rssi_decibel", "HiLink signal RSSI in decibel")
	metricSignalRSRP           = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_signal_rsrp_decibel", "HiLink signal RSRP in decibel")
	metricSignalRSRQ           = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_signal_rsrq_decibel", "HiLink signal RSRQ in decibel")
	metricSignalSINR           = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_signal_sinr_decibel", "HiLink signal SINR in decibel")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricBalance.With("serial_number", sn).Describe(ch)
	metricLimitInternetTraffic.With("serial_number", sn).Describe(ch)
	metricSignalRSSI.With("serial_number", sn).Describe(ch)
	metricSignalRSRP.With("serial_number", sn).Describe(ch)
	metricSignalRSRQ.With("serial_number", sn).Describe(ch)
	metricSignalSINR.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.SerialNumber()
	if sn == "" {
		return
	}

	metricBalance.With("serial_number", sn).Collect(ch)
	metricLimitInternetTraffic.With("serial_number", sn).Collect(ch)
	metricSignalRSSI.With("serial_number", sn).Collect(ch)
	metricSignalRSRP.With("serial_number", sn).Collect(ch)
	metricSignalRSRQ.With("serial_number", sn).Collect(ch)
	metricSignalSINR.With("serial_number", sn).Collect(ch)
}
