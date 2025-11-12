package daemon

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/aid297/aid/filesystem/filesystemV2"
	"github.com/aid297/aid/operation/operationV2"
)

// Daemon 守护进程服务提供者
type Daemon struct{}

// Launch 启动守护进程
func (*Daemon) Launch(title, logDir, logFilename string) {
	var (
		err  error
		dir  *filesystemV2.Dir  = filesystemV2.APP.Dir.NewByAbs(logDir)
		file *filesystemV2.File = filesystemV2.APP.File.NewByAbs(logDir, operationV2.NewTernary(operationV2.TrueValue(logFilename), operationV2.FalseValue("daemon.log")).GetByValue(logFilename != ""))
		fp   *os.File
	)

	if syscall.Getppid() == 1 {
		if err := os.Chdir("./"); err != nil {
			panic(err)
		}
		// syscall.Umask(0) // TODO TEST
		return
	}

	if !dir.GetExist() {
		if err = dir.Create(os.ModePerm).Error(); err != nil {
			log.Fatalf("【启动失败】创建日志目录失败：%s", err.Error())
		}
	}

	if fp, err = os.OpenFile(file.GetFullPath(), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644); err != nil {
		log.Fatalf("【启动失败】创建总日志失败：%s", err.Error())
	}
	defer func() { _ = fp.Close() }()
	cmd := exec.Command(os.Args[0], os.Args[1:]...)
	// cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true} // TODO TEST
	cmd.SysProcAttr = &syscall.SysProcAttr{} // TODO TEST
	cmd.Stdout = fp
	cmd.Stderr = fp
	cmd.Stdin = nil
	if err = cmd.Start(); err != nil {
		log.Fatalf("【启动失败】%s", err.Error())
	}

	if _, err = fmt.Fprintf(
		fp,
		"--------------------------------------------------\r\n%s 程序启动成功 [进程号->%d] 启动于：%s\r\n",
		title,
		cmd.Process.Pid,
		time.Now().Format(string(time.DateTime+".000")),
	); err != nil {
		log.Fatalf("【启动失败】写入日志失败：%s", err.Error())
	}

	os.Exit(0)
}
