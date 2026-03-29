package video

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type VideoInfo struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	FrameRate  string `json:"r_frame_rate"`
	FrameCount string `json:"nb_frames"`
}

// Compare paths to ensure same metadata
func SameVideoInfo(refPath, disPath string) (bool, error) {
	// Get video information
	refInfo, err := GetVideoInfo(refPath)
	if err != nil {
		return false, err
	}

	disInfo, err := GetVideoInfo(disPath)
	if err != nil {
		return false, err
	}

	// Compare resolutions
	if refInfo.Width != disInfo.Width || refInfo.Height != disInfo.Height {
		return false, fmt.Errorf(
			"reference and distorted files have different resolutions: %dx%d, %dx%d",
			refInfo.Width,
			refInfo.Height,
			disInfo.Width,
			disInfo.Height,
		)
	}

	// Compare frame rate strings
	if refInfo.FrameRate != disInfo.FrameRate {
		return false, fmt.Errorf(
			"reference and distorted files have different framerates: %s, %s",
			refInfo.FrameRate,
			disInfo.FrameRate,
		)
	}

	// Compare frame count
	if refInfo.FrameCount != disInfo.FrameCount {
		return false, fmt.Errorf(
			"reference and distorted files have different frame counts: %s, %s",
			refInfo.FrameCount,
			disInfo.FrameCount,
		)
	}

	return true, nil
}

// Get information of video
func GetVideoInfo(path string) (VideoInfo, error) {
	// Local struct to hold ffprobe output
	type ffprobeOut struct {
		Streams []VideoInfo `json:"streams"`
	}

	// Get json formatted video information
	cmd := exec.Command(
		"ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height,r_frame_rate,nb_frames",
		"-of", "json",
		path,
	)

	// Get command output
	out, err := cmd.Output()
	if err != nil {
		return VideoInfo{}, err
	}

	// Unmarshal to struct
	var res ffprobeOut
	if err := json.Unmarshal(out, &res); err != nil {
		return VideoInfo{}, err
	}

	// Ensure only one video stream exists
	if len(res.Streams) != 1 {
		return VideoInfo{}, fmt.Errorf("only one video stream permitted in file: %s", path)
	}

	return res.Streams[0], nil
}
