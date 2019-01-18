package storage

import (
	"github.com/kihamo/shadow"
)

type Component interface {
	shadow.Component

	NamespacePath(namespace string) (path string, err error)
	SaveURLToFile(namespace, url string, force bool) (string, error)
}
