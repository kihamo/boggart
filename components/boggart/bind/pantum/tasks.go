package pantum

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kihamo/boggart/providers/pantum/client/om"
	"github.com/kihamo/go-workers"
)

func (b *Bind) Tasks() []workers.Task {
	taskUpdater := b.Workers().WrapTaskIsOnline(b.taskUpdater)
	taskUpdater.SetTimeout(b.config.UpdaterTimeout)
	taskUpdater.SetRepeats(-1)
	taskUpdater.SetRepeatInterval(b.config.UpdaterInterval)
	taskUpdater.SetName("updater")

	return []workers.Task{
		taskUpdater,
	}
}

func (b *Bind) taskUpdater(ctx context.Context) error {
	response, err := b.provider.Om.GetDatabase(om.NewGetDatabaseParamsWithContext(ctx))
	if err != nil {
		return err
	}

	if b.Meta().SerialNumber() == "" {
		for _, p := range response.Payload {
			switch p.Name {
			case "omSerialNumber":
				b.Meta().SetSerialNumber(p.Value)
			case "omMACAddress":
				b.Meta().SetMACAsString(p.Value)
			}
		}
	}

	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	for _, p := range response.Payload {
		switch p.Name {
		case "omTonerRemain":
			if v, e := strconv.ParseUint(p.Value, 10, 64); e == nil {
				metricTonerRemain.With("serial_number", sn).Set(float64(v))

				if e := b.MQTT().PublishAsync(ctx, b.config.TopicTonerRemain.Format(sn), v); e != nil {
					err = fmt.Errorf("send mqtt message about toner remain failed: %w", err)
				}
			} else {
				err = fmt.Errorf("oarse toner remain value failed: %w", err)
			}

		case "omPrinterStatus":
			//0: "Initialization"
			//1: "Sleep"
			//2: "Warming Up…"
			//3: "Ready"
			//4: "Printing …"
			//5: "Error"
			//134: "Canceling …"

			var v string

			switch p.Value {
			case "0":
				v = "initialization"
			case "1":
				v = "sleep"
			case "2":
				v = "warming up"
			case "3":
				v = "ready"
			case "4":
				v = "printing"
			case "5":
				v = "error"
			case "134":
				v = "canceling"
			default:
				v = "unknown"
			}

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicPrinterStatus.Format(sn), v); e != nil {
				err = fmt.Errorf("send mqtt message about printer status failed: %w", err)
			}

		case "omCartridgeStatus":
			//0: "Normal"
			//1: "Cartridge not detected"
			//2: "Cartridge mismatch"
			//3: "Cartridge Life Expired"
			//4: "Toner Low"
			//5: "Unknown Status"

			var v string

			switch p.Value {
			case "0":
				v = "normal"
			case "1":
				v = "not detected"
			case "2":
				v = "mismatch"
			case "3":
				v = "expired"
			case "4":
				v = "toner low"
			default:
				v = "unknown"
			}

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicCartridgeStatus.Format(sn), v); e != nil {
				err = fmt.Errorf("send mqtt message about cartridge status failed: %w", err)
			}

		case "omDrumStatus":
			//0: "Normal"
			//1: "Drum unit is uninstalled"
			//2: "Type of the drum unit is unmatched"
			//3: "Drum unit is invalid"
			//4: "Life of the drum unit will come to an end"
			//5: "Unknown Status"

			var v string

			switch p.Value {
			case "0":
				v = "normal"
			case "1":
				v = "uninstalled"
			case "2":
				v = "unmatched"
			case "3":
				v = "invalid"
			case "4":
				v = "expired"
			default:
				v = "unknown"
			}

			if e := b.MQTT().PublishAsync(ctx, b.config.TopicDrumStatus.Format(sn), v); e != nil {
				err = fmt.Errorf("send mqtt message about drum status failed: %w", err)
			}
		}
	}

	return nil
}
