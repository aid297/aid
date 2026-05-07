//go:build unix || linux || darwin || freebsd || netbsd || openbsd

package daemon

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

// Launch 启动守护进程
func (my *Daemon) Launch() {
	var (
		err error
		fp  = my.bootLogFile()
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

	if err = my.afterLaunch(fp, cmd.Process.Pid); err != nil {
		log.Fatalf("【启动失败】%s", err.Error())
	}

	os.Exit(0)
}
