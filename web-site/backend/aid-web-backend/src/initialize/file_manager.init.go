package initialize

import (
	"github.com/aid297/aid/filesystem/filesystemV4"
	"github.com/aid297/aid/web-site/backend/aid-web-backend/src/global"
)

type FileManagerInitialize struct{}

func (FileManagerInitialize) Boot() {
	if dir := filesystemV4.NewDir(filesystemV4.Rel(global.CONFIG.FileManager.Dir)); !dir.GetExist() {
		dir.Create(filesystemV4.Flag(0644))
	}
}
