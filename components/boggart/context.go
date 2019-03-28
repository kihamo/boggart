package boggart

import (
	"context"
)

type contextKey string

var (
	i18nDomain = contextKey("i18n-domain")
)

func ContextWithI18nDomain(ctx context.Context, domain string) context.Context {
	return context.WithValue(ctx, i18nDomain, domain)
}

func I18nDomainFromContext(c context.Context) string {
	v := c.Value(i18nDomain)
	if v != nil {
		return v.(string)
	}

	return ""
}
