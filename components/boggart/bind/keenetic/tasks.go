package keenetic

import (
	"context"
	"fmt"
	"time"

	"github.com/kihamo/boggart/components/boggart/tasks"
	"github.com/kihamo/boggart/components/mqtt"
	"github.com/kihamo/boggart/providers/keenetic/client/show"
	"go.uber.org/multierr"
)

const (
	TaskNameSerialNumber = "serial-number"
	TaskNameUpdater      = "updater"
	TaskNameHotspotSync  = "connections"
)

func (b *Bind) Tasks() []tasks.Task {
	return []tasks.Task{
		tasks.NewTask().
			WithName(TaskNameSerialNumber).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskSerialNumberHandler),
				),
			).
			WithSchedule(
				tasks.ScheduleWithSuccessLimit(
					tasks.ScheduleWithDuration(tasks.ScheduleNow(), time.Second*30),
					1,
				),
			),
	}
}

func (b *Bind) taskSerialNumberHandler(ctx context.Context) error {
	defaults, err := b.client.Show.ShowDefaults(show.NewShowDefaultsParamsWithContext(ctx))
	if err != nil {
		return fmt.Errorf("get defaults value failed: %w", err)
	}

	b.Meta().SetSerialNumber(defaults.Payload.Serial)

	cfg := b.config()

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName(TaskNameUpdater).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerFuncFromShortToLong(b.taskUpdaterHandler),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.UpdaterInterval)),
	)
	if err != nil {
		return fmt.Errorf("register task "+TaskNameUpdater+" failed: %w", err)
	}

	_, err = b.Workers().RegisterTask(
		tasks.NewTask().
			WithName(TaskNameHotspotSync).
			WithHandler(
				b.Workers().WrapTaskHandlerIsOnline(
					tasks.HandlerWithTimeout(
						tasks.HandlerFunc(b.taskHotspotSyncHandler),
						cfg.ReadinessTimeout,
					),
				),
			).
			WithSchedule(tasks.ScheduleWithDuration(tasks.ScheduleNow(), cfg.HotspotSyncInterval)),
	)
	if err != nil {
		return fmt.Errorf("register task "+TaskNameHotspotSync+" failed: %w", err)
	}

	if cfg.TopicSyslog != "" {
		err = b.MQTT().Subscribe(mqtt.NewSubscriber(cfg.TopicSyslog, 0, b.callbackMQTTSyslog))

		if err != nil {
			return fmt.Errorf("register subscriber for syslog failed: %w", err)
		}
	}

	return nil
}

func (b *Bind) taskUpdaterHandler(ctx context.Context) (err error) {
	sn := b.Meta().SerialNumber()

	system, err := b.client.Show.ShowSystem(show.NewShowSystemParamsWithContext(ctx))
	if err == nil {
		metricUpTime.With("serial_number", sn).Set(float64(system.Payload.Uptime))
		metricCPULoad.With("serial_number", sn).Set(float64(system.Payload.Cpuload))
		metricMemoryAvailable.With("serial_number", sn).Set(float64(system.Payload.Memfree * 1024))
		metricMemoryUsage.With("serial_number", sn).Set(float64(system.Payload.Memtotal-system.Payload.Memfree) * 1024)
	} else {
		err = multierr.Append(err, fmt.Errorf("get system info failed: %w", err))
	}

	return err
}

func (b *Bind) taskHotspotSyncHandler(ctx context.Context, meta tasks.Meta, _ tasks.Task) (err error) {

	fmt.Println("Sync start")

	hotspot, err := b.client.Show.ShowIPHotspot(show.NewShowIPHotspotParamsWithContext(ctx))
	if err != nil {
		return fmt.Errorf("get IP hostpot failed: %w", err)
	}

	var (
		wifiClients int
		item        *storeItem
	)

	currentVersion := meta.Attempts()
	sn := b.Meta().SerialNumber()
	cfg := b.config()

	/*
		FIXME: есть проблема, что после реального отключения устройства и сообщения об этом в syslog
		само api по хотспоту продолжает некоторое время (примерно 1 минута) отдавать активность подключения,
		но это только для НЕ зарегистрированных устройств, зарегистрированные будут висеть всегда.
		Можно по полю link определять жив ли коннект, как и делает интерфейс. Итого
		Active - зеленая кнопочка в интерфейсе кинетика
		Link - реальное подключение в интерфейсе кинетика
	*/
	for _, host := range hotspot.Payload.Host {
		item = &storeItem{
			version: currentVersion,
			host:    host,
		}

		b.hotspotConnections.Store(item.ID(), item)

		if host.Active && host.Ssid != "" {
			wifiClients++
		}
	}

	metricWifiClients.With("serial_number", sn).Set(float64(wifiClients))

	// Теперь рассылаем актуальные версии в MQTT, а то, что меньше текущей версии еще и удаляем
	b.hotspotConnections.Range(func(key, value interface{}) bool {
		si := value.(*storeItem)

		// старая версия, которую надо удалять, чтобы список не пух
		if si.version < currentVersion {
			b.hotspotConnections.Delete(key)

			// фактически удаляем запись в MQTT
			_ = b.MQTT().PublishAsyncWithoutCache(ctx, cfg.TopicHotspotState.Format(sn, si.ID()), nil)
		} else {
			// постим актуальный статус в MQTT
			_ = b.MQTT().PublishAsync(ctx, cfg.TopicHotspotState.Format(sn, si.ID()), si)
		}

		return true
	})

	// подписываемся на собственные топики, это нужно для того чтобы при первом запуске
	// проверить старые записи в mqtt, где вероятна ситуация, когда соединение отвалилось
	// и новый бинд его не видит никак так как апи не отдает а словари пусты, но он продолжает
	// висеть в mqtt с прошлым статусом
	b.hotspotZombieKiller.Do(func() {
		err = b.MQTT().Subscribe(
			mqtt.NewSubscriber(
				cfg.TopicHotspotState.Format(sn),
				0,
				b.callbackMQTTHotspotSearchZombies,
			),
		)

	})

	if err != nil {
		err = fmt.Errorf("register mqtt subscriber for search zombies of hotspot failed: %w", err)
		b.hotspotZombieKiller.Reset()
	}

	return err
}
