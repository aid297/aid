package initialize

import (
	"github.com/aid297/aid/v2/filesystem"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/global"
)

type FileManagerInitialize struct{}

func (FileManagerInitialize) Boot() {
	if dir := filesystem.NewDir(filesystem.Rel(global.CONFIG.FileManager.Dir)); !dir.GetExist() {
		dir.Create(filesystem.Flag(0644))
	}
}
