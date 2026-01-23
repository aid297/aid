package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/aid297/aid/filesystem/filesystemV4"
	"github.com/aid297/aid/operation/operationV2"
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

// GetDaemonOnce 获取单例
func GetDaemonOnce() *Daemon { o.Do(func() { i = &Daemon{} }); return i }

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

// Launch 启动守护进程
func (my *Daemon) Launch() {
	var (
		err  error
		dir  filesystemV4.Filesystemer
		file filesystemV4.Filesystemer
		fp   *os.File
	)

	if my.enable && my.dir != "" {
		dir = filesystemV4.NewDir(filesystemV4.Rel(my.dir))
		file = filesystemV4.NewFile(filesystemV4.Abs(operationV2.NewTernary(operationV2.TrueValue(my.filename), operationV2.FalseValue("daemon.log")).GetByValue(my.filename != "")))
	}

	if syscall.Getppid() == 1 {
		if err = os.Chdir("./"); err != nil {
			panic(err)
		}
		// syscall.Umask(0)
		return
	}

	if dir != nil && !dir.GetExist() {
		if err = dir.Create(filesystemV4.Mode(os.ModePerm)).GetError(); err != nil {
			log.Fatalf("【启动失败】创建日志目录失败：%s", err.Error())
		}
	}

	if file != nil {
		if fp, err = os.OpenFile(file.GetFullPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			log.Fatalf("【启动失败】创建总日志失败：%s", err.Error())
		}
		defer func() { _ = fp.Close() }()
	}
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{} // TODO TEST
	if fp != nil {
		cmd.Stdout = fp
		cmd.Stderr = fp
	}
	cmd.Stdin = nil
	if err = cmd.Start(); err != nil {
		log.Fatalf("【启动失败】%s", err.Error())
	}

	if fp != nil {
		if _, err = fmt.Fprintf(
			fp,
			"--------------------------------------------------\r\n%s 程序启动成功 [进程号->%d] 启动于：%s\r\n",
			my.title,
			cmd.Process.Pid,
			time.Now().Format(time.DateTime+".000"),
		); err != nil {
			log.Fatalf("【启动失败】写入日志失败：%s", err.Error())
		}
	}

	os.Exit(0)
}
