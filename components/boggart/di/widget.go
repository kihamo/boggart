package di

import (
	"context"
	"net/url"
	"sync"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

type WidgetHandler interface {
	WidgetHandler(*dashboard.Response, *dashboard.Request)
	WidgetAssetFS() *assetfs.AssetFS
}

type WidgetContainerSupport interface {
	SetWidget(*WidgetContainer)
	Widget() *WidgetContainer
}

func WidgetContainerBind(bind boggart.Bind) (*WidgetContainer, bool) {
	if support, ok := bind.(WidgetContainerSupport); ok {
		container := support.Widget()
		return container, container != nil
	}

	return nil, false
}

type WidgetBind struct {
	mutex     sync.RWMutex
	container *WidgetContainer
}

func (b *WidgetBind) SetWidget(container *WidgetContainer) {
	b.mutex.Lock()
	b.container = container
	b.mutex.Unlock()
}

func (b *WidgetBind) Widget() *WidgetContainer {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.container
}

type WidgetContainer struct {
	dashboard.Handler
	bind boggart.BindItem
}

func NewWidgetContainer(bind boggart.BindItem) *WidgetContainer {
	return &WidgetContainer{
		bind: bind,
	}
}

func (c *WidgetContainer) Bind() boggart.Bind {
	return c.bind.Bind()
}

func (c *WidgetContainer) Handle(w *dashboard.Response, r *dashboard.Request) {
	if h, ok := c.Bind().(WidgetHandler); ok {
		r = r.WithContext(dashboard.ContextWithTemplateNamespace(r.Context(), c.TemplateNamespace()))
		r = r.WithContext(boggart.ContextWithI18nDomain(r.Context(), c.I18nDomain()))

		h.WidgetHandler(w, r)
	}
}

func (c *WidgetContainer) AssetFS() *assetfs.AssetFS {
	if w, ok := c.Bind().(WidgetHandler); ok {
		return w.WidgetAssetFS()
	}

	return nil
}

func (c *WidgetContainer) TemplateNamespace() string {
	return boggart.ComponentName + "-bind-" + c.bind.ID()
}

func (c *WidgetContainer) I18nDomain() string {
	return boggart.ComponentName + "-bind-" + c.bind.ID()
}

func (c *WidgetContainer) Translate(ctx context.Context, messageID string, context string, format ...interface{}) string {
	return i18n.Locale(ctx).Translate(boggart.I18nDomainFromContext(ctx), messageID, context, format...)
}

func (c *WidgetContainer) TranslatePlural(ctx context.Context, singleID, pluralID string, number int, context string, format ...interface{}) string {
	return i18n.Locale(ctx).TranslatePlural(boggart.I18nDomainFromContext(ctx), singleID, pluralID, number, context, format...)
}

func (c *WidgetContainer) URL() (*url.URL, error) {
	return nil, nil
}

func (c *WidgetContainer) FlashError(r *dashboard.Request, messageID string, context string, format ...interface{}) {
	r.Session().FlashBag().Error(c.Translate(r.Context(), messageID, context, format...))
}

func (c *WidgetContainer) FlashSuccess(r *dashboard.Request, messageID string, context string, format ...interface{}) {
	r.Session().FlashBag().Success(c.Translate(r.Context(), messageID, context, format...))
}

func (c *WidgetContainer) FlashInfo(r *dashboard.Request, messageID string, context string, format ...interface{}) {
	r.Session().FlashBag().Info(c.Translate(r.Context(), messageID, context, format...))
}
