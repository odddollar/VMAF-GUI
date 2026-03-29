package video

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os/exec"
	"time"
)

// Contains progress information from ffmpeg output
type Progress struct {
	Frame   int
	FPS     int
	Time    time.Duration
	Speed   float64
	Elapsed time.Duration
}

// Run vmaf calculation with progress updates
func RunVMAF(ctx context.Context, ref, dist string) (<-chan Progress, <-chan error, error) {
	// Create channel to push progress status through
	progressChan := make(chan Progress)
	errChan := make(chan error, 1)

	// Create ffmpeg command to run vmaf calculation
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-hide_banner",
		"-loglevel", "error",
		"-stats",
		"-i", dist,
		"-i", ref,
		"-lavfi", "libvmaf=n_threads=8:log_path=vmaf.json:log_fmt=json",
		"-f", "null", "-",
	)

	// Command will output over stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	// Start ffmpeg
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	go func() {
		defer close(progressChan)
		defer close(errChan)
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
	}()

	return progressChan, errChan, nil
}
