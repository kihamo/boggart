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
	metricSignalLevel          = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_signal_level", "HiLink signal level")
	metricTotalConnectTime     = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_total_connection_time_seconds", "HiLink total connection time in seconds")
	metricTotalDownload        = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_total_download_bytes", "HiLink total download in bytes")
	metricTotalUpload          = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_total_upload_bytes", "HiLink total upload in bytes")
	metricSMSUnread            = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_sms_unread", "HiLink SMS unread")
	metricSMSInbox             = snitch.NewGauge(boggart.ComponentName+"_bind_hilink_sms_inbox", "HiLink SMS inbox")
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
	metricSignalLevel.With("serial_number", sn).Describe(ch)
	metricTotalConnectTime.With("serial_number", sn).Describe(ch)
	metricTotalDownload.With("serial_number", sn).Describe(ch)
	metricTotalUpload.With("serial_number", sn).Describe(ch)
	metricSMSUnread.With("serial_number", sn).Describe(ch)
	metricSMSInbox.With("serial_number", sn).Describe(ch)
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
	metricSignalLevel.With("serial_number", sn).Collect(ch)
	metricTotalConnectTime.With("serial_number", sn).Collect(ch)
	metricTotalDownload.With("serial_number", sn).Collect(ch)
	metricTotalUpload.With("serial_number", sn).Collect(ch)
	metricSMSUnread.With("serial_number", sn).Collect(ch)
	metricSMSInbox.With("serial_number", sn).Collect(ch)
}
