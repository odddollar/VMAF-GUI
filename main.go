package main

import (
	"VMAF-GUI/ui"
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Progress struct {
	Frame   int
	FPS     float64
	Time    time.Duration
	Speed   float64
	Elapsed time.Duration
}

func parseProgress(line string) (Progress, bool) {
	var p Progress

	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return p, false
	}

	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		if !strings.Contains(t, "=") {
			continue
		}

		kv := strings.SplitN(t, "=", 2)
		key := kv[0]
		val := kv[1]

		// Handle "key=" followed by value in next token
		if val == "" && i+1 < len(tokens) {
			val = tokens[i+1]
			i++ // Consume next token
		}

		val = strings.TrimSpace(val)

		switch key {
		case "frame":
			if v, err := strconv.Atoi(val); err == nil {
				p.Frame = v
			}
		case "fps":
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				p.FPS = v
			}
		case "time":
			if d, err := parseTime(val); err == nil {
				p.Time = d
			}
		case "speed":
			val = strings.TrimSuffix(val, "x")
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				p.Speed = v
			}
		case "elapsed":
			if d, err := parseTime(val); err == nil {
				p.Elapsed = d
			}
		}
	}

	return p, true
}

func parseTime(s string) (time.Duration, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, strconv.ErrSyntax
	}

	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	sec, _ := strconv.ParseFloat(parts[2], 64)

	d := time.Duration(h)*time.Hour +
		time.Duration(m)*time.Minute +
		time.Duration(sec*float64(time.Second))

	return d, nil
}

func RunVMAF(ref, dist string) (<-chan Progress, error) {
	out := make(chan Progress)

	cmd := exec.Command(
		"ffmpeg",
		"-hide_banner",
		"-loglevel", "error",
		"-stats",
		"-i", dist,
		"-i", ref,
		"-lavfi", "libvmaf=n_threads=8:log_path=vmaf.json:log_fmt=json",
		"-f", "null", "-",
	)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		defer close(out)

		r := bufio.NewReader(stderr)
		var buf bytes.Buffer

		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			if b == '\r' || b == '\n' {
				line := buf.String()
				if p, ok := parseProgress(line); ok {
					out <- p
				}
				buf.Reset()
			}

			buf.WriteByte(b)
		}

		cmd.Wait()
	}()

	return out, nil
}

func main() {
	// ch, err := RunVMAF("reference.mp4", "distorted.mp4")
	// if err != nil {
	// 	panic(err)
	// }

	// for line := range ch {
	// 	fmt.Println(line)
	// }

	// Create and run UI
	u := ui.Ui{}
	u.NewUI()
	u.Run()
}
