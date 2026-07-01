//go:build unix || linux || darwin || freebsd || netbsd || openbsd

package daemons

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// Launch 启动守护进程
func (*Daemon) Launch() {
	var (
		err error
		fp  = i.bootLogFile()
		cmd *exec.Cmd
	)
	defer func() {
		if fp != nil {
			_ = fp.Close()
		}
	}()

	if syscall.Getppid() == 1 {
		if err = os.Chdir("./"); err != nil {
			panic(err)
		}
		return
	}

	cmd = exec.Command(os.Args[0], os.Args[1:]...)
	{
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		cmd.Stdin = nil
		if fp != nil {
			cmd.Stdout = fp
			cmd.Stderr = fp
		}

		if err = cmd.Start(); err != nil {
			log.Fatalf("【启动失败】%s", err.Error())
		}
	}

	if err = i.afterLaunch(fp, cmd.Process.Pid); err != nil {
		log.Fatalf("【启动失败】%s", err.Error())
	}

	log.Println("【启动成功】Daemon 模式")

	os.Exit(0)
}

// LaunchWithCondition 通过条件决定是否启动守护进程
func (*Daemon) LaunchWithCondition(condition bool, bootFn func()) {
	if bootFn == nil {
		return
	}

	if condition {
		i.Launch()
	} else {
		bootFn()
	}
}
