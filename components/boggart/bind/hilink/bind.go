package hilink

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/boggart/components/boggart/providers/hilink"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/client/ussd"
	"github.com/kihamo/boggart/components/boggart/providers/hilink/models"
)

type Bind struct {
	boggart.BindBase
	boggart.BindMQTT

	config        *Config
	client        *hilink.Client
	balanceRegexp *regexp.Regexp
}

func (b *Bind) USSD(ctx context.Context, content string) (string, error) {
	if content == "" {
		return "", nil
	}

	params := ussd.NewSendUSSDParamsWithContext(ctx).
		WithRequest(&models.USSD{
			Content: content,
		})

	_, err := b.client.Ussd.SendUSSD(params)
	if err != nil {
		return "", err
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()

		default:
			response, err := b.client.Ussd.GetUSSD(ussd.NewGetUSSDParamsWithContext(ctx))
			if err == nil && response.Payload.Content != "" {
				return response.Payload.Content, nil
			}

			time.Sleep(time.Second)
		}
	}

	return "", err
}

func (b *Bind) Balance(ctx context.Context) (float64, error) {
	content, err := b.USSD(ctx, b.config.BalanceUSSD)
	if err != nil {
		return -1, err
	}

	match := b.balanceRegexp.FindStringSubmatch(content)
	for i, name := range b.balanceRegexp.SubexpNames() {
		if name == "balance" {
			return strconv.ParseFloat(match[i], 64)
		}
	}

	return 0, errors.New("balance not found")
}
