package hilink

import (
	"context"
	"errors"

	"github.com/kihamo/boggart/components/boggart/installer"
	"github.com/kihamo/boggart/components/boggart/installer/openhab"
)

func (b *Bind) InstallersSupport() []installer.System {
	return []installer.System{
		installer.SystemOpenHab,
	}
}

func (b *Bind) InstallerSteps(_ context.Context, s installer.System) ([]installer.Step, error) {
	meta := b.Meta()
	sn := meta.SerialNumber()
	if sn == "" {
		return nil, errors.New("serial number is empty")
	}

	itemPrefix := openhab.ItemPrefixFromBindMeta(meta)
	cfg := b.config()
	channels := make([]*openhab.Channel, 0, 17)

	const (
		idSignalLevel                     = "SignalLevel"
		idOperator                        = "Operator"
		idSignalRSSI                      = "SignalRSSI"
		idSignalRSRP                      = "SignalRSRP"
		idSignalRSRQ                      = "SignalRSRQ"
		idSignalSINR                      = "SignalSINR"
		idMobileCurrentConnectionDuration = "MobileCurrentConnectionDuration"
		idMobileCurrentDownload           = "MobileCurrentDownload"
		idMobileCurrentUpload             = "MobileCurrentUpload"
		idMobileTotalConnectionDuration   = "MobileTotalConnectionDuration"
		idMobileTotalDownload             = "MobileTotalDownload"
		idMobileTotalUpload               = "MobileTotalUpload"
		idBalance                         = "Balance"
		idSMSInbox                        = "SMSInbox"
		idSMSUnread                       = "SMSUnread"
		idSMSLast                         = "SMSLast"
		idSMSLastContent                  = "SMSLast_Content"
		idSMSLastPhone                    = "SMSLast_Phone"
		idSMSLastDate                     = "SMSLast_Date"
		idLimitInternetTraffic            = "LimitInternetTraffic"
	)

	transformHumanSeconds := openhab.StepDefaultTransformHumanSeconds.Base()
	transformHumanBytes := openhab.StepDefaultTransformHumanBytes.Base()

	// system-updater

	channels = append(channels,
		openhab.NewChannel(idSignalLevel, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicSignalLevel.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idSignalLevel, openhab.ItemTypeNumber).
					WithLabel("Signal level").
					WithIcon("qualityofservice"),
			),
		openhab.NewChannel(idMobileCurrentConnectionDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicMobileCurrentConnectionDuration.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMobileCurrentConnectionDuration, openhab.ItemTypeNumber).
					WithLabel("Mobile current connection time [JS("+transformHumanSeconds+"):%s]").
					WithIcon("time"),
			),
		openhab.NewChannel(idMobileCurrentDownload, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicMobileCurrentDownload.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMobileCurrentDownload, openhab.ItemTypeNumber).
					WithLabel("Mobile current download [JS("+transformHumanBytes+"):%s]").
					WithIcon("returnpipe"),
			),
		openhab.NewChannel(idMobileCurrentUpload, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicMobileCurrentUpload.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMobileCurrentUpload, openhab.ItemTypeNumber).
					WithLabel("Mobile current upload [JS("+transformHumanBytes+"):%s]").
					WithIcon("flowpipe"),
			),
		openhab.NewChannel(idMobileTotalConnectionDuration, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicMobileTotalConnectionDuration.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMobileTotalConnectionDuration, openhab.ItemTypeNumber).
					WithLabel("Mobile total connection time [JS("+transformHumanSeconds+"):%s]").
					WithIcon("time"),
			),
		openhab.NewChannel(idMobileTotalDownload, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicMobileTotalDownload.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMobileTotalDownload, openhab.ItemTypeNumber).
					WithLabel("Mobile total download [JS("+transformHumanBytes+"):%s]").
					WithIcon("returnpipe"),
			),
		openhab.NewChannel(idMobileTotalUpload, openhab.ChannelTypeNumber).
			WithStateTopic(cfg.TopicMobileTotalUpload.Format(sn)).
			AddItems(
				openhab.NewItem(itemPrefix+idMobileTotalUpload, openhab.ItemTypeNumber).
					WithLabel("Mobile total upload [JS("+transformHumanBytes+"):%s]").
					WithIcon("flowpipe"),
			),
	)

	if !b.operator.IsEmpty() {
		channels = append(channels,
			openhab.NewChannel(idOperator, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicOperator.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idOperator, openhab.ItemTypeString).
						WithLabel("Operator"),
				),
		)
	}

	if b.simStatus.Load() == 1 {
		channels = append(channels,
			openhab.NewChannel(idSignalRSSI, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicSignalRSSI.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idSignalRSSI, openhab.ItemTypeNumber).
						WithLabel("Signal RSSI [%d dBm]").
						WithIcon("qualityofservice"),
				),
			openhab.NewChannel(idSignalRSRP, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicSignalRSRP.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idSignalRSRP, openhab.ItemTypeNumber).
						WithLabel("Signal RSRP [%d dBm]").
						WithIcon("qualityofservice"),
				),
			openhab.NewChannel(idSignalRSRQ, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicSignalRSRQ.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idSignalRSRQ, openhab.ItemTypeNumber).
						WithLabel("Signal RSRQ [%d dB]").
						WithIcon("qualityofservice"),
				),
			openhab.NewChannel(idSignalSINR, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicSignalSINR.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idSignalSINR, openhab.ItemTypeNumber).
						WithLabel("Signal SINR [%d dB]").
						WithIcon("qualityofservice"),
				),
		)
	}

	// balance-updater
	if _, _, err := b.Workers().TaskInfoByName("balance-updater"); err == nil {
		channels = append(channels,
			openhab.NewChannel(idBalance, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicBalance.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idBalance, openhab.ItemTypeNumber).
						WithLabel("Balance [%.2f ₽]").
						WithIcon("price"),
				),
		)
	}

	// sms-checker
	if _, _, err := b.Workers().TaskInfoByName("sms-checker"); err == nil {
		channels = append(channels,
			openhab.NewChannel(idSMSInbox, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicSMSInbox.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idSMSInbox, openhab.ItemTypeNumber).
						WithLabel("SMS inbox"),
				),
			openhab.NewChannel(idSMSUnread, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicSMSUnread.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idSMSUnread, openhab.ItemTypeNumber).
						WithLabel("SMS unread"),
				),
			openhab.NewChannel(idSMSLast, openhab.ChannelTypeString).
				WithStateTopic(cfg.TopicSMSLast.Format(sn)).
				AddItems(
					// {"Content":"text","Date":"2022-05-21 23:05:46","Index":40062,"Phone":"Tele2","SaveType":4,"SmsType":2}
					openhab.NewItem(itemPrefix+idSMSLastContent, openhab.ItemTypeString).
						WithLabel("Last SMS content [JSONPATH($.Content):%s]"),
					openhab.NewItem(itemPrefix+idSMSLastPhone, openhab.ItemTypeString).
						WithLabel("Last SMS phone [JSONPATH($.Phone):%s]"),
					openhab.NewItem(itemPrefix+idSMSLastDate, openhab.ItemTypeString).
						WithLabel("Last SMS date [JSONPATH($.Date):%s]"). // TODO: to datetime type
						WithIcon("calendar"),
				),
			openhab.NewChannel(idLimitInternetTraffic, openhab.ChannelTypeNumber).
				WithStateTopic(cfg.TopicLimitInternetTraffic.Format(sn)).
				AddItems(
					openhab.NewItem(itemPrefix+idLimitInternetTraffic, openhab.ItemTypeNumber).
						WithLabel("Internet traffic limit [JS("+transformHumanBytes+"):%s]").
						WithIcon("line"),
				),
		)
	}

	return openhab.StepsByBind(b, []installer.Step{
		openhab.StepDefault(openhab.StepDefaultTransformHumanBytes),
		openhab.StepDefault(openhab.StepDefaultTransformHumanSeconds),
	}, channels...)
}
