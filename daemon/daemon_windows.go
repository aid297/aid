//go:build windows

package daemon

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

const _DETACHED_PROCESS = 0x00000008 // Windows API 常量

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

	cmd = exec.Command(os.Args[0], os.Args[1:]...)
	{
		cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: _DETACHED_PROCESS, HideWindow: true}
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
