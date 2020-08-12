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
	WidgetHandler(w *dashboard.Response, r *dashboard.Request)
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

func (b *WidgetContainer) Bind() boggart.Bind {
	return b.bind.Bind()
}

func (b *WidgetContainer) Handle(w *dashboard.Response, r *dashboard.Request) {
	if h, ok := b.Bind().(WidgetHandler); ok {
		r = r.WithContext(dashboard.ContextWithTemplateNamespace(r.Context(), b.TemplateNamespace()))
		r = r.WithContext(boggart.ContextWithI18nDomain(r.Context(), b.I18nDomain()))

		h.WidgetHandler(w, r)
	}
}

func (b *WidgetContainer) AssetFS() *assetfs.AssetFS {
	if w, ok := b.Bind().(WidgetHandler); ok {
		return w.WidgetAssetFS()
	}

	return nil
}

func (b *WidgetContainer) TemplateNamespace() string {
	return boggart.ComponentName + "-bind-" + b.bind.ID()
}

func (b *WidgetContainer) I18nDomain() string {
	return boggart.ComponentName + "-bind-" + b.bind.ID()
}

func (b *WidgetContainer) Translate(ctx context.Context, messageID string, context string, format ...interface{}) string {
	return i18n.Locale(ctx).Translate(boggart.I18nDomainFromContext(ctx), messageID, context, format...)
}

func (b *WidgetContainer) TranslatePlural(ctx context.Context, singleID, pluralID string, number int, context string, format ...interface{}) string {
	return i18n.Locale(ctx).TranslatePlural(boggart.I18nDomainFromContext(ctx), singleID, pluralID, number, context, format...)
}

func (b *WidgetContainer) URL() (*url.URL, error) {
	return nil, nil
}
