package initialize

import (
	"github.com/aid297/aid/v2/filesystems"
	"github.com/aid297/aid/v2/web-site/backend/aid-web-backend/src/global"
)

type FileManagerInitialize struct{}

func (FileManagerInitialize) Boot() {
	if dir := filesystems.NewDir(filesystems.Rel(global.CONFIG.FileManager.Dir)); !dir.GetExist() {
		dir.Create(filesystems.Flag(0644))
	}
}
