package video

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"io"
	"os/exec"
)

// Given two video files get the Nth frame from both (0-indexed)
func GetFramePair(
	ctx context.Context,
	refPath, disPath string,
	refInfo VideoInfo,
	frameIndex int,
) (*image.NRGBA, *image.NRGBA, error) {
	// Filter options to normalise and output raw rgb frames
	filter := fmt.Sprintf(
		"[0:v]settb=AVTB,setpts=PTS-STARTPTS,fps=%s,scale=%d:%d:flags=bicubic,format=rgb24,select=eq(n\\,%d)[dis];"+
			"[1:v]settb=AVTB,setpts=PTS-STARTPTS,fps=%s,scale=%d:%d:flags=bicubic,format=rgb24,select=eq(n\\,%d)[ref];"+
			"[dis][ref]concat=n=2:v=1:a=0[out]",
		refInfo.FrameRate, refInfo.Width, refInfo.Height, frameIndex,
		refInfo.FrameRate, refInfo.Width, refInfo.Height, frameIndex,
	)

	// Run for both videos simultaneously and output over pipes
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg",
		"-hide_banner",
		"-loglevel", "error",
		"-i", disPath,
		"-i", refPath,
		"-filter_complex", filter,
		"-map", "[out]",
		"-frames:v", "2",
		"-vsync", "0",
		"-f", "rawvideo",
		"-pix_fmt", "rgb24",
		"pipe:1",
	)

	// Will receive both frames as single buffer
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	// Read command errors
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, err
	}

	// Read stderr concurrently
	var stderrBuf []byte
	done := make(chan struct{})
	go func() {
		stderrBuf, _ = io.ReadAll(stderr)
		close(done)
	}()

	// Start command
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	// Buffer to store output frames
	frameSize := refInfo.Width * refInfo.Height * 3
	totalSize := frameSize * 2
	buf := make([]byte, totalSize)

	// Read all stdout to buffer
	_, err = io.ReadFull(stdout, buf)
	if err != nil {
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		<-done
		return nil, nil, fmt.Errorf("%v: %s", err, string(stderrBuf))
	}

	// Wait for command completion
	if err := cmd.Wait(); err != nil {
		<-done
		if ctx.Err() != nil {
			return nil, nil, ctx.Err()
		}
		return nil, nil, fmt.Errorf("%v: %s", err, string(stderrBuf))
	}

	// Ensure stderr fully read
	<-done

	// Split buffers
	disBuf := buf[:frameSize]
	refBuf := buf[frameSize:]

	// Process byte arrays to actual frames
	disImg := rgbToNRGBA(disBuf, refInfo.Width, refInfo.Height)
	refImg := rgbToNRGBA(refBuf, refInfo.Width, refInfo.Height)

	// Return reference frame first
	return refImg, disImg, nil
}

// Convert byte array to image
func rgbToNRGBA(buf []byte, width, height int) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	i := 0
	for y := range height {
		for x := range width {
			img.SetNRGBA(x, y, color.NRGBA{
				R: buf[i],
				G: buf[i+1],
				B: buf[i+2],
				A: 255,
			})
			i += 3
		}
	}
	return img
}
