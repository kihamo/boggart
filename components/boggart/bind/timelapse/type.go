package timelapse

import (
	"os"
	"path/filepath"

	"github.com/kihamo/boggart/components/boggart"
)

type Type struct{}

func (t Type) CreateBind() boggart.Bind {
	return &Bind{}
}

func (t Type) DashboardTemplateFunctions() map[string]interface{} {
	return map[string]interface{}{
		"file_id": FileID,
	}
}

func FileID(file os.FileInfo) string {
	return filepath.Join(file.ModTime().Format(SubDirectoryNameLayout), file.Name())
}
