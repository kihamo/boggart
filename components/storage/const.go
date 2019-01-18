package storage

import (
	"errors"
	"os"
)

const (
	ComponentName    = "storage"
	ComponentVersion = "1.0.0"

	RouteFileStoragePrefix = "/" + ComponentName + "/file/"
)

var (
	Separator = string(os.PathSeparator)

	ErrFileAlreadyExist = errors.New("file already exist")
)
