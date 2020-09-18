package openweathermap

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

var (
	rateLimitersLock = new(sync.Mutex)
	rateLimiters     = make(map[string]*Limiter, 1)
)

const (
	PriceUnlimited    = 0
	PriceFree         = 60
	PriceStartup      = 600
	PriceDeveloper    = 3000
	PriceProfessional = 30000
	PriceEnterprise   = 200000

	LimiterInterval = time.Minute
)

type Limiter struct {
	rateLimiter *rate.Limiter
}

func NewLimiter(apiKey string, price int) (l *Limiter) {
	var d rate.Limit

	switch price {
	case PriceFree, PriceStartup, PriceDeveloper, PriceProfessional, PriceEnterprise:
		d = rate.Every(LimiterInterval)
	default:
		d = rate.Inf
	}

	rateLimitersLock.Lock()
	defer rateLimitersLock.Unlock()

	l, ok := rateLimiters[apiKey]
	if !ok {
		l = &Limiter{
			rateLimiter: rate.NewLimiter(d, price),
		}

		rateLimiters[apiKey] = l
	} else {
		// увеличенный рейтлимит всегда переопределяет весь лимитатор в целом, более широкое условие
		if l.rateLimiter.Burst() < price {
			l.rateLimiter.SetBurst(price)
		}

		if l.rateLimiter.Limit() < d {
			l.rateLimiter.SetLimit(d)
		}
	}

	return l
}

func (l *Limiter) LimitInterval() time.Duration {
	if l.rateLimiter.Limit() != rate.Inf {
		return LimiterInterval
	}

	return 0
}

func (l *Limiter) Wait(ctx context.Context) error {
	return l.rateLimiter.Wait(ctx)
}
