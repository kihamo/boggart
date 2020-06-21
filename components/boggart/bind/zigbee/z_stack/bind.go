package z_stack

import (
	"context"
	"strconv"
	"sync"

	"github.com/kihamo/boggart/atomic"
	"github.com/kihamo/boggart/components/boggart/di"
	"github.com/kihamo/boggart/protocols/connection"
	"github.com/kihamo/boggart/providers/zigbee/z_stack"
)

type Bind struct {
	di.MetaBind
	di.MQTTBind
	di.LoggerBind
	di.ProbesBind
	di.WorkersBind

	config *Config

	disconnected *atomic.BoolNull
	client       *z_stack.Client
	onceClient   *atomic.Once
	lock         sync.RWMutex
	done         chan struct{}
}

func (b *Bind) getClient(ctx context.Context) (_ *z_stack.Client, err error) {
	b.onceClient.Do(func() {
		b.disconnected.Nil()
		var conn connection.Conn

		conn, err = connection.New(b.config.ConnectionDSN)
		if err != nil {
			b.disconnected.True()
		}

		b.client = z_stack.New(conn)
		err = b.client.Boot(ctx)
	})

	if err != nil {
		b.onceClient.Reset()
	}

	return b.client, err
}

func (b *Bind) Run() error {
	b.onceClient.Reset()
	doneCh := make(chan struct{})

	b.lock.Lock()
	b.done = doneCh
	b.lock.Unlock()

	client, err := b.getClient(context.Background())
	if err != nil {
		return err
	}

	go func() {
		defer b.disconnected.True()

		watcher := client.Watch()

		for {
			select {
			case frame := <-watcher.NextFrame():
				b.Logger().Debug("Received Zigbee message",
					"length", frame.Length(),
					"type", frame.Type(),
					"sub-system", frame.SubSystem(),
					"command-id", frame.CommandID(),
					"data", frame.Data(),
					"fcs", frame.FCS(),
				)

				if frame.CommandID() == z_stack.CommandAfIncomingMessage {
					message, err := client.AfIncomingMessage(frame)
					if err != nil {
						b.Logger().Warn("Parse received message", "error", err.Error())
						continue
					}

					ctx := context.Background()
					sourceAddress := strconv.FormatUint(uint64(message.SrcAddr), 10)
					sn := b.Meta().SerialNumber()

					if sn != "" {
						metricLinkQuality.With("serial_number", sn).With("srcaddr", sourceAddress).Set(float64(message.LinkQuality))
						_ = b.MQTT().Publish(ctx, b.config.TopicLinkQuality.Format(sn, sourceAddress), message.LinkQuality)
					}

					if message.Frame.Payload.Report != nil && len(*message.Frame.Payload.Report) > 0 {
						report := (*message.Frame.Payload.Report)[0]

						switch report.AttributeID {
						// battery
						case 65282:
							elements := report.AttributeData.(z_stack.TypeStruct).Elements
							if len(elements) > 1 {
								if sn != "" {
									voltage := elements[1].Value.(uint16)
									percent := ToPercentageCR2032(voltage)

									metricBatteryPercent.With("serial_number", sn).With("srcaddr", sourceAddress).Set(float64(percent))
									_ = b.MQTT().Publish(ctx, b.config.TopicBatteryPercent.Format(sn, sourceAddress), percent)

									metricBatteryVoltage.With("serial_number", sn).With("srcaddr", sourceAddress).Set(float64(voltage))
									_ = b.MQTT().Publish(ctx, b.config.TopicBatteryVoltage.Format(sn, sourceAddress), voltage)
								}
							} else {
								b.Logger().Warn("Battery element not found", "elements", len(elements))
							}

							// double click
						case 32768:
							_ = b.MQTT().PublishWithoutCache(ctx, b.config.TopicOnOff.Format(sn, sourceAddress), true)
							_ = b.MQTT().Publish(ctx, b.config.TopicClick.Format(sn, sourceAddress), report.AttributeData)

							// hmmm
						default:
							if sn != "" && report.DataType == z_stack.DataTypeBoolean {
								value := report.AttributeData.(bool)

								_ = b.MQTT().PublishWithoutCache(ctx, b.config.TopicOnOff.Format(sn, sourceAddress), value)

								if value {
									_ = b.MQTT().Publish(ctx, b.config.TopicClick.Format(sn, sourceAddress), 1)
								}
							}
						}
					}
				}

			case err := <-watcher.NextError():
				b.Logger().Warn("Watcher error", "error", err.Error())

			case <-watcher.Done():
				return

			case <-doneCh:
				return
			}
		}
	}()

	return nil
}

func (b *Bind) Close() (err error) {
	b.lock.RLock()
	defer b.lock.RUnlock()

	if b.client != nil {
		// close(b.done)
		return b.client.Close()
	}

	return nil
}

func ToPercentageCR2032(voltage uint16) (percentage uint16) {
	switch _ = voltage; {
	case voltage < 2100:
		percentage = 0

	case voltage < 2440:
		percentage = 6 - ((2440-voltage)*6)/340

	case voltage < 2740:
		percentage = 18 - ((2740-voltage)*12)/300

	case voltage < 2900:
		percentage = 42 - ((2900-voltage)*24)/160

	case voltage < 3000:
		percentage = 100 - ((3000-voltage)*58)/100

	case voltage >= 3000:
		percentage = 100
	}

	return percentage
}