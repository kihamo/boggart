package xmeye

import (
	"context"
	"sync/atomic"
	"time"
)

type AlertStreaming struct {
	client         *Client
	ctx            context.Context
	tickerInterval time.Duration
	alerts         chan *AlarmInfo
	errors         chan error
}

func (c *Client) AlarmStart(ctx context.Context) error {
	_, err := c.Call(ctx, CmdGuardRequest, nil)

	if err == nil {
		atomic.StoreUint32(&c.alarmStarted, 1)
	}

	return err
}

func (c *Client) AlarmStop(ctx context.Context) error {
	_, err := c.Call(ctx, CmdUnGuardRequest, nil)

	if err == nil {
		atomic.StoreUint32(&c.alarmStarted, 0)
	}

	return err
}

func (c *Client) AlarmInfo(ctx context.Context) (*AlarmInfo, error) {
	var result struct {
		Response
		AlarmInfo AlarmInfo
	}

	err := c.CallWithResult(ctx, CmdAlarmRequest, nil, &result)

	if err != nil {
		return nil, err
	}

	if result.AlarmInfo.Event == "" {
		return nil, nil
	}

	return &result.AlarmInfo, nil
}

func (c *Client) AlarmStreaming(ctx context.Context, interval time.Duration) *AlertStreaming {
	s := &AlertStreaming{
		client:         c,
		ctx:            ctx,
		tickerInterval: interval,
		alerts:         make(chan *AlarmInfo),
		errors:         make(chan error),
	}
	go s.start()

	return s
}

func (s *AlertStreaming) NextAlarm() <-chan *AlarmInfo {
	return s.alerts
}

func (s *AlertStreaming) NextError() <-chan error {
	return s.errors
}

func (s *AlertStreaming) start() {
	ticker := time.NewTicker(s.tickerInterval)

	defer func() {
		ticker.Stop()
		close(s.alerts)
		close(s.errors)
	}()

	var (
		alarmEnable bool
		err         error
	)

	for {
		select {
		case <-s.ctx.Done():
			return

		case <-ticker.C:
			if !alarmEnable {
				if err = s.client.AlarmStart(s.ctx); err != nil {
					s.errors <- err
					continue
				}

				alarmEnable = true
			}

			alert, err := s.client.AlarmInfo(s.ctx)
			if err != nil {
				s.errors <- err
			} else if alert != nil {
				s.alerts <- alert
			}
		}
	}
}
