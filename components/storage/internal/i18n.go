package internal

import (
	"io"

	"github.com/kihamo/shadow/components/i18n"
)

func (c *Component) I18n() map[string][]io.ReadSeeker {
	fs := c.AssetFS()
	fs.Prefix = "locales"

	return i18n.FromAssetFS(fs)
}
