package di

import (
	"context"
	"net/url"
	"strings"
	"sync"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/kihamo/boggart/components/boggart"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
	"github.com/kihamo/shadow/components/i18n"
)

type contextKey string

var (
	contextKeyWidget = contextKey("widget")
)

func ContextWithWidget(ctx context.Context, widget *WidgetContainer) context.Context {
	return context.WithValue(ctx, contextKeyWidget, widget)
}

func WidgetFromContext(c context.Context) *WidgetContainer {
	if v := c.Value(contextKeyWidget); v != nil {
		return v.(*WidgetContainer)
	}

	return nil
}

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

	bindItem  boggart.BindItem
	configApp config.Component
}

func NewWidgetContainer(bindItem boggart.BindItem, configApp config.Component) *WidgetContainer {
	return &WidgetContainer{
		bindItem:  bindItem,
		configApp: configApp,
	}
}

func (c *WidgetContainer) globalVariables() map[string]interface{} {
	vars := map[string]interface{}{
		"Config": struct{}{},
		"Meta":   struct{}{},
	}

	if bindSupport, ok := ConfigContainerBind(c.bindItem.Bind()); ok {
		vars["Config"] = bindSupport.Bind()
	}

	if bindSupport, ok := MetaContainerBind(c.bindItem.Bind()); ok {
		vars["Meta"] = bindSupport
	}

	return vars
}

func (c *WidgetContainer) Bind() boggart.Bind {
	return c.bindItem.Bind()
}

func (c *WidgetContainer) HandleAllowed() bool {
	status := c.bindItem.Status()

	return status.IsStatusOnline() || status.IsStatusOffline()
}

func (c *WidgetContainer) Handle(w *dashboard.Response, r *dashboard.Request) {
	if !c.HandleAllowed() {
		c.NotFound(w, r)
		return
	}

	if h, ok := c.Bind().(WidgetHandler); ok {
		r = r.WithContext(ContextWithWidget(r.Context(), c))
		r = r.WithContext(boggart.ContextWithI18nDomain(r.Context(), c.I18nDomain()))
		r = r.WithContext(dashboard.ContextWithTemplateNamespace(r.Context(), c.TemplateNamespace()))
		r = r.WithContext(dashboard.ContextWithRequest(r.Context(), r))

		h.WidgetHandler(w, r)
	}
}

func (c *WidgetContainer) Render(ctx context.Context, view string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{}, 1)
	}

	data["Bind"] = c.globalVariables()

	c.Handler.Render(ctx, view, data)
}

func (c *WidgetContainer) RenderLayout(ctx context.Context, view, layout string, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{}, 1)
	}

	data["Bind"] = c.globalVariables()

	c.Handler.RenderLayout(ctx, view, layout, data)
}

func (c *WidgetContainer) AssetFS() *assetfs.AssetFS {
	if w, ok := c.Bind().(WidgetHandler); ok {
		return w.WidgetAssetFS()
	}

	return nil
}

func (c *WidgetContainer) TemplateNamespace() string {
	return boggart.ComponentName + "-bind-" + c.bindItem.ID()
}

func (c *WidgetContainer) I18nDomain() string {
	return boggart.ComponentName + "-bind-" + c.bindItem.ID()
}

func (c *WidgetContainer) Translate(ctx context.Context, messageID string, context string, format ...interface{}) string {
	return i18n.Locale(ctx).Translate(boggart.I18nDomainFromContext(ctx), messageID, context, format...)
}

func (c *WidgetContainer) TranslatePlural(ctx context.Context, singleID, pluralID string, number int, context string, format ...interface{}) string {
	return i18n.Locale(ctx).TranslatePlural(boggart.I18nDomainFromContext(ctx), singleID, pluralID, number, context, format...)
}

func (c *WidgetContainer) URL(vs map[string]string) (*url.URL, error) {
	u, err := url.Parse(c.configApp.String(boggart.ConfigExternalURL))
	if err != nil {
		return nil, err
	}

	u.Path = "/" + boggart.ComponentName + "/widget/" + c.bindItem.ID()

	values := u.Query()

	if keysConfig := c.configApp.String(boggart.ConfigAccessKeys); keysConfig != "" {
		if keys := strings.Split(keysConfig, ","); len(keys) > 0 {
			values.Add(boggart.AccessKeyName, keys[0])
		}
	}

	for k, v := range vs {
		values.Add(k, v)
	}

	u.RawQuery = values.Encode()

	return u, nil
}

func (c *WidgetContainer) FlashError(r *dashboard.Request, messageID interface{}, context string, format ...interface{}) {
	var id string

	if e, ok := messageID.(error); ok {
		id = e.Error()
	} else {
		id = messageID.(string)
	}

	r.Session().FlashBag().Error(c.Translate(r.Context(), id, context, format...))
}

func (c *WidgetContainer) FlashSuccess(r *dashboard.Request, messageID string, context string, format ...interface{}) {
	r.Session().FlashBag().Success(c.Translate(r.Context(), messageID, context, format...))
}

func (c *WidgetContainer) FlashInfo(r *dashboard.Request, messageID string, context string, format ...interface{}) {
	r.Session().FlashBag().Info(c.Translate(r.Context(), messageID, context, format...))
}
