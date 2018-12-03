package internal

import (
	"github.com/elazarl/go-bindata-assetfs"
)

func (c *Component) AssetFS() *assetfs.AssetFS {
	return assetFS()
}
