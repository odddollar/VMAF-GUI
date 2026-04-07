package video

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
)

// Run vmaf calculation with progress updates
func RunVMAF(ctx context.Context, refPath, disPath string) (<-chan Progress, <-chan error, <-chan struct{}, error) {
	// Create channel to push progress status through
	progressChan := make(chan Progress)
	errChan := make(chan error, 1)
	doneChan := make(chan struct{})

	// Get reference video info
	refInfo, err := GetVideoInfo(refPath)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create ffmpeg filter
	filter := fmt.Sprintf(
		"[0:v]settb=AVTB,setpts=PTS-STARTPTS,fps=%s,scale=%d:%d:flags=bicubic,format=%s[dis];"+
			"[1:v]settb=AVTB,setpts=PTS-STARTPTS,fps=%s,format=%s[ref];"+
			"[dis][ref]libvmaf=n_threads=8:log_path=vmaf.json:log_fmt=json",
		refInfo.FrameRate,
		refInfo.Width, refInfo.Height,
		refInfo.PixFmt,
		refInfo.FrameRate,
		refInfo.PixFmt,
	)

	// Create ffmpeg command to run vmaf calculation
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-hide_banner",
		"-loglevel", "error",
		"-stats",
		"-i", disPath,
		"-i", refPath,
		"-lavfi", filter,
		"-f", "null", "-",
	)

	// Command will output over stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	// Start ffmpeg
	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}

	// Parse progress
	go func() {
		defer close(progressChan)
		var sentErr bool // Ensures at most one error sent

		// Read from command output
		r := bufio.NewReader(stderr)
		var buf bytes.Buffer

		for {
			// Read byte by byte until empty
			b, err := r.ReadByte()
			if err != nil {
				// Stop if context was cancelled
				if ctx.Err() != nil {
					break
				}

				// Send read errors that aren't no more data
				if err != io.EOF {
					sentErr = true
					errChan <- err
				}
				break
			}

			// Reached end of line
			if b == '\r' || b == '\n' {
				// Get whole line
				line := buf.String()

				// Send data over channel
				p, err := parseProgress(line)
				if err == nil { // Ignore parsing errors as they're not fatal
					progressChan <- p
				}

				// Reset and read next line
				buf.Reset()
				continue
			}

			// Not end of line so write to buffer
			buf.WriteByte(b)
		}

		// Wait for command to finish
		if err := cmd.Wait(); err != nil {
			// Don't report manual cancellation
			if ctx.Err() != nil || sentErr {
				return
			}

			errChan <- err
			return
		}

		// Report when finished sucessfully
		close(doneChan)
	}()

	return progressChan, errChan, doneChan, nil
}
