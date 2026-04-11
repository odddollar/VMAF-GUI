package video

import (
	"os/exec"
	"runtime"
	"syscall"
)

// Hide command prompt in Windows
func hideCmdWindow(cmd *exec.Cmd) {
	if runtime.GOOS == "windows" {
		cmd.SysProcAttr = &syscall.SysProcAttr{
			HideWindow: true,
		}
	}
}
