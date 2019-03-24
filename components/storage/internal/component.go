package internal

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/kihamo/boggart/components/storage"
	"github.com/kihamo/shadow"
	"github.com/kihamo/shadow/components/config"
	"github.com/kihamo/shadow/components/dashboard"
)

var (
	hasher = md5.New()
)

type Component struct {
	config config.Component
	routes []dashboard.Route
}

func (c *Component) Name() string {
	return storage.ComponentName
}

func (c *Component) Version() string {
	return storage.ComponentVersion
}

func (c *Component) Dependencies() []shadow.Dependency {
	return []shadow.Dependency{
		{
			Name:     config.ComponentName,
			Required: true,
		},
	}
}

func (c *Component) Init(a shadow.Application) error {
	c.config = a.GetComponent(config.ComponentName).(config.Component)
	return nil
}

func (c *Component) Run(a shadow.Application, ready chan<- struct{}) error {
	return nil
}

func (c *Component) NamespacePath(namespace string) (path string, err error) {
	namespaces := c.config.String(storage.ConfigFileNameSpaces)
	if namespaces == "" {
		return "", errors.New("namespace not found")
	}

	for _, n := range strings.Split(namespaces, ",") {
		parts := strings.Split(n, ":")
		if len(parts) < 2 {
			continue
		}

		ns := strings.ToLower(strings.TrimSpace(parts[0]))
		if ns != namespace {
			continue
		}

		path = strings.TrimRight(strings.TrimSpace(parts[1]), storage.Separator)
		path = filepath.FromSlash(path)
		path, err = filepath.Abs(path)

		if err == nil {
			break
		}
	}

	if path == "" {
		err = errors.New("namespace not found")
	}

	return path, err
}

func (c *Component) SaveURLToFile(namespace, url string, force bool) (string, error) {
	// save
	path, err := c.NamespacePath(namespace)
	if err != nil {
		return "", nil
	}

	if _, err := io.WriteString(hasher, url); err != nil {
		return "", err
	}

	// TODO: extension of file
	id := hex.EncodeToString(hasher.Sum(nil))
	filePath := filepath.Join(path, id)

	file, err := os.Open(filePath)
	if err != nil {
		defer file.Close()
	} else {
		if force {
			if err := os.Remove(filePath); err != nil {
				return "", err
			}
		} else {
			return storage.RouteFileStoragePrefix + namespace + "/" + id, storage.ErrFileAlreadyExist
		}
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// load
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	_, err = io.Copy(out, response.Body)
	if err != nil {
		return "", err
	}

	return storage.RouteFileStoragePrefix + namespace + "/" + id, nil
}
