package video

import (
	"bytes"
	"os/exec"
	"strings"
)

// Check if ffmpeg command is available
func FFmpegAvailable() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// Check if vmaf available in ffmpeg
func VMAFAvailable() bool {
	cmd := exec.Command("ffmpeg", "-filters")

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return false
	}

	// Look for libvmaf filter in output
	return strings.Contains(out.String(), "libvmaf")
}
