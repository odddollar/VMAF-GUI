package video

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Parse progress output string from ffmpeg
func parseProgress(line string) (Progress, error) {
	var p Progress

	// Split by whitespace
	// Whitespace output for ffmpeg changes by number of characters in value
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return p, fmt.Errorf("no tokens in line")
	}

	// Iterate through split tokens
	for i := 0; i < len(tokens); i++ {
		t := tokens[i]

		// Skip non-key tokens
		if !strings.Contains(t, "=") {
			continue
		}

		// Create key/value pair
		kv := strings.SplitN(t, "=", 2)
		key := kv[0]
		val := kv[1]

		// Handle "key=" followed by value in next token
		if val == "" && i+1 < len(tokens) {
			val = tokens[i+1]
			i++ // Consume next token
		}

		// Set value of output struct
		switch key {
		case "frame":
			// Get current frame
			v, err := strconv.Atoi(val)
			if err != nil {
				return Progress{}, err
			}
			p.Frame = v
		case "fps":
			// Get current frames per second
			v, err := strconv.Atoi(val)
			if err != nil {
				return Progress{}, err
			}
			p.FPS = v
		case "time":
			// Get current time of video
			v, err := parseTime(val)
			if err != nil {
				return Progress{}, err
			}
			p.Time = v
		case "speed":
			// Get current processing speed
			v, err := strconv.ParseFloat(strings.TrimSuffix(val, "x"), 64)
			if err != nil {
				return Progress{}, err
			}
			p.Speed = v
		case "elapsed":
			// Get elapsed processing time
			v, err := parseTime(val)
			if err != nil {
				return Progress{}, err
			}
			p.Elapsed = v
		}
	}

	return p, nil
}

// Parse string into time format
func parseTime(s string) (time.Duration, error) {
	// Split string
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return 0, fmt.Errorf("unable to parse time: %s", s)
	}

	// Parse to Go's Duration type
	return time.ParseDuration(parts[0] + "h" + parts[1] + "m" + parts[2] + "s")
}

// Parse fps "X/X" to float
func parseFPS(rate string) (float64, error) {
	var num, den float64
	_, err := fmt.Sscanf(rate, "%f/%f", &num, &den)
	if err != nil || den == 0 {
		return 0, fmt.Errorf("unable to parse frame rate: %s", rate)
	}
	return num / den, nil
}
