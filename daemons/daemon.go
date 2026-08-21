package daemons

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"

	"github.com/aid297/aid/v3/filesystems"
	"github.com/aid297/aid/v3/operations"
)

// Daemon 守护进程服务提供者
type Daemon struct {
	title    string
	dir      string
	filename string
	enable   bool
}

var (
	o sync.Once
	i *Daemon
)

// OnceDaemon 获取单例
func OnceDaemon() *Daemon { o.Do(func() { i = &Daemon{} }); return i }

// SetTitle 设置标题
func (*Daemon) SetTitle(title string) *Daemon { i.title = title; return i }

// SetLogDir 设置日志目录
func (*Daemon) SetLogDir(logDir string) *Daemon { i.dir = logDir; return i }

// SetLogFilename 设置日志文件名
func (*Daemon) SetLogFilename(logFilename string) *Daemon { i.filename = logFilename; return i }

// SetLog 设置日志
func (*Daemon) SetLog(dir, filename string) *Daemon {
	i.dir = dir
	i.filename = filename
	i.enable = true
	return i
}

// SetLogEnable 设置日志开关
func (*Daemon) SetLogEnable(enable bool) *Daemon { i.enable = enable; return i }

// bootLogFile 启动日志
func (*Daemon) bootLogFile() (fp *os.File) {
	var (
		err  error
		dir  filesystems.Filesystem
		file filesystems.Filesystem
	)

	if i.enable && i.dir != "" {
		dir = filesystems.NewDir(filesystems.Rel(i.dir))
		file = filesystems.NewFile(
			filesystems.Abs(
				path.Join(
					dir.GetFullPath(),
					operations.NewTernary(operations.TrueValue(i.filename), operations.FalseValue("daemons.log")).GetByValue(i.filename != ""),
				),
			),
		)
	}

	if dir != nil && !dir.GetExist() {
		if err = dir.Create(filesystems.Mode(os.ModePerm)).GetError(); err != nil {
			log.Fatalf("【启动失败】创建日志目录失败：%s", err.Error())
		}
	}

	if file != nil {
		if fp, err = os.OpenFile(file.GetFullPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			log.Fatalf("【启动失败】创建总日志失败：%s", err.Error())
		}
	}

	return
}

// afterLaunch 成功启动守护进程后
func (*Daemon) afterLaunch(fp *os.File, pid int) (err error) {
	if fp != nil {
		if _, err = fmt.Fprintf(
			fp,
			"--------------------------------------------------\r\n%s 程序启动成功 [进程号->%d] 启动于：%s\r\n",
			i.title,
			pid,
			time.Now().Format(time.DateTime+".000"),
		); err != nil {
			log.Fatalf("【启动失败】写入日志失败：%s", err.Error())
		}
	}

	return
}
