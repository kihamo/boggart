package pantum

import (
	"context"
	"fmt"
	"strconv"

	"github.com/kihamo/boggart/components/boggart/tasks"
	client "github.com/kihamo/boggart/providers/pantum"
	"github.com/kihamo/boggart/providers/pantum/client/om"
)

func (b *Bind) Tasks() []tasks.Task {
	cfg := b.config()

	return []tasks.Task{
		tasks.NewTask().
			WithName("updater").
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
						cfg.UpdaterTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	}
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) error {
	response, err := b.provider.Om.GetDatabase(om.NewGetDatabaseParamsWithContext(ctx))
	if err != nil {
		return err
	}

	database := client.DatabaseToMap(response.Payload)

	if b.Meta().SerialNumber() == "" {
		if value, ok := database["omSerialNumber"]; ok {
			b.Meta().SetSerialNumber(value.(string))
		}

		if value, ok := database["omMACAddress"]; ok {
			b.Meta().SetMACAsString(value.(string))
		}
	}

	sn := b.Meta().SerialNumber()
	if sn == "" {
		return nil
	}

	cfg := b.config()

	var product client.ProductID
	if value, ok := database["omProductID"]; ok {
		product = client.ProductIDConvert(value.(string))

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicProductID.Format(sn), product.String()); e != nil {
			err = fmt.Errorf("send mqtt message about printer status failed: %w", err)
		}
	}

	if flag, ok := database["omErrorFlag"]; ok {
		if modules, ok := database["omStatusModule"]; ok {
			status := client.PrinterStatus(flag.(string), modules.(string))

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicPrinterStatus.Format(sn), status); e != nil {
				err = fmt.Errorf("send mqtt message about printer status failed: %w", err)
			}
		}
	}

	if value, ok := database["omTonerRemain"]; ok {
		if v, e := strconv.ParseUint(value.(string), 10, 64); e == nil {
			metricTonerRemain.With("serial_number", sn).Set(float64(v))

			if e := b.MQTT().PublishAsync(ctx, cfg.TopicTonerRemain.Format(sn), v); e != nil {
				err = fmt.Errorf("send mqtt message about toner remain failed: %w", err)
			}
		} else {
			err = fmt.Errorf("oarse toner remain value failed: %w", err)
		}
	}

	if value, ok := database["omCartridgeStatus"]; ok {
		status := client.CartridgeStatusConvert(value.(string))

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicCartridgeStatus.Format(sn), status.String()); e != nil {
			err = fmt.Errorf("send mqtt message about cartridge status failed: %w", err)
		}
	}

	if value, ok := database["omDrumStatus"]; ok {
		status := client.DrumStatusConvert(value.(string))

		if e := b.MQTT().PublishAsync(ctx, cfg.TopicDrumStatus.Format(sn), status.String()); e != nil {
			err = fmt.Errorf("send mqtt message about drum status failed: %w", err)
		}
	}

	return err
}
