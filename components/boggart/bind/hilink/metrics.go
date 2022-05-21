package hilink

import (
	"github.com/kihamo/snitch"
)

var (
	metricBalance                    = snitch.NewGauge("balance_rubles", "Balance in rubles")
	metricLimitInternetTraffic       = snitch.NewGauge("limit_internet_traffic_bytes", "Limit internet traffic in bytes")
	metricSignalRSSI                 = snitch.NewGauge("signal_rssi_decibel", "Signal RSSI in decibel")
	metricSignalRSRP                 = snitch.NewGauge("signal_rsrp_decibel", "Signal RSRP in decibel")
	metricSignalRSRQ                 = snitch.NewGauge("signal_rsrq_decibel", "Signal RSRQ in decibel")
	metricSignalSINR                 = snitch.NewGauge("signal_sinr_decibel", "Signal SINR in decibel")
	metricSignalLevel                = snitch.NewGauge("signal_level", "Signal level")
	metricMobileTotalConnectDuration = snitch.NewGauge("mobile_total_connection_duration_seconds", "Total mobile connection duration in seconds")
	metricMobileTotalDownload        = snitch.NewGauge("mobile_total_download_bytes", "Total mobile download in bytes")
	metricMobileTotalUpload          = snitch.NewGauge("mobile_total_upload_bytes", "Total mobile upload in bytes")
	metricSMSUnread                  = snitch.NewGauge("sms_unread", "SMS unread")
	metricSMSInbox                   = snitch.NewGauge("sms_inbox", "SMS inbox")
)

func (b *Bind) Describe(ch chan<- *snitch.Description) {
	sn := b.Meta().SerialNumber()
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
	metricMobileTotalConnectDuration.With("serial_number", sn).Describe(ch)
	metricMobileTotalDownload.With("serial_number", sn).Describe(ch)
	metricMobileTotalUpload.With("serial_number", sn).Describe(ch)
	metricSMSUnread.With("serial_number", sn).Describe(ch)
	metricSMSInbox.With("serial_number", sn).Describe(ch)
}

func (b *Bind) Collect(ch chan<- snitch.Metric) {
	sn := b.Meta().SerialNumber()
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
	metricMobileTotalConnectDuration.With("serial_number", sn).Collect(ch)
	metricMobileTotalDownload.With("serial_number", sn).Collect(ch)
	metricMobileTotalUpload.With("serial_number", sn).Collect(ch)
	metricSMSUnread.With("serial_number", sn).Collect(ch)
	metricSMSInbox.With("serial_number", sn).Collect(ch)
}
