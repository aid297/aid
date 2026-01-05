package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/aid297/aid/filesystem/filesystemV3"
	"github.com/aid297/aid/operation/operationV2"
)

// Daemon 守护进程服务提供者
type Daemon struct {
	title       string
	logDir      string
	logFilename string
	logEnable   bool
}

var (
	daemonOnce sync.Once
	daemonIns  *Daemon
)

// Once 获取单例
func (*Daemon) Once() *Daemon {
	daemonOnce.Do(func() { daemonIns = &Daemon{} })
	return daemonIns
}

// SetTitle 设置标题
func (*Daemon) SetTitle(title string) *Daemon {
	daemonIns.title = title
	return daemonIns
}

// SetLogDir 设置日志目录
func (*Daemon) SetLogDir(logDir string) *Daemon {
	daemonIns.logDir = logDir
	return daemonIns
}

// SetLogFilename 设置日志文件名
func (*Daemon) SetLogFilename(logFilename string) *Daemon {
	daemonIns.logFilename = logFilename
	return daemonIns
}

// SetLog 设置日志
func (*Daemon) SetLog(dir, filename string) *Daemon {
	daemonIns.logDir = dir
	daemonIns.logFilename = filename
	daemonIns.logEnable = true
	return daemonIns
}

// SetLogEnable 设置日志开关
func (*Daemon) SetLogEnable(enable bool) *Daemon {
	daemonIns.logEnable = enable
	return daemonIns
}

// Launch 启动守护进程
func (my *Daemon) Launch() {
	var (
		err  error
		dir  *filesystemV3.Dir
		file *filesystemV3.File
		fp   *os.File
	)

	if my.logEnable && my.logDir != "" {
		dir = filesystemV3.APP.Dir.New(filesystemV3.APP.DirAttr.IsRel.SetRel(), filesystemV3.APP.DirAttr.Path.Set(my.logDir))
		file = filesystemV3.APP.File.New(filesystemV3.APP.FileAttr.IsRel.SetAbs(), filesystemV3.APP.FileAttr.Path.Set(operationV2.NewTernary(operationV2.TrueValue(my.logFilename), operationV2.FalseValue("daemon.log")).GetByValue(my.logFilename != "")))
	}

	if syscall.Getppid() == 1 {
		if err = os.Chdir("./"); err != nil {
			panic(err)
		}
		// syscall.Umask(0)
		return
	}

	if dir != nil && !dir.Exist {
		if err = dir.Create(filesystemV3.DirMode(os.ModePerm)).Error; err != nil {
			log.Fatalf("【启动失败】创建日志目录失败：%s", err.Error())
		}
		// if err = dir.Create(os.ModePerm).Error(); err != nil {
		// 	log.Fatalf("【启动失败】创建日志目录失败：%s", err.Error())
		// }
	}

	if file != nil {
		if fp, err = os.OpenFile(file.FullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
			log.Fatalf("【启动失败】创建总日志失败：%s", err.Error())
		}
		// if fp, err = os.OpenFile(file.GetFullPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
		// 	log.Fatalf("【启动失败】创建总日志失败：%s", err.Error())
		// }
		defer func() { _ = fp.Close() }()
	}
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	// cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // TODO TEST
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
			time.Now().Format(string(time.DateTime+".000")),
		); err != nil {
			log.Fatalf("【启动失败】写入日志失败：%s", err.Error())
		}
	}

	os.Exit(0)
}
