package video

import (
	"os/exec"
	"strings"
)

// Check if given command is available
func CommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

// Check if vmaf available in ffmpeg
func VMAFAvailable() bool {
	cmd := exec.Command("ffmpeg", "-filters")

	// Get command output
	out, err := cmd.Output()
	if err != nil {
		return false
	}

	// Look for libvmaf filter in output
	return strings.Contains(string(out), "libvmaf")
}
