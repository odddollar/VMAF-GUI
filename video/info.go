package video

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type VideoInfo struct {
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	FrameRate  string `json:"r_frame_rate"`
	FrameCount string `json:"nb_frames"`
	PixFmt     string `json:"pix_fmt"`
}

// Compare paths to ensure same metadata
// Must be same frame rate and frame count as frames must align for comparisons
// Frame rate will be normalised to reference's to enforce CFR
// Resolution and pixel format also normalised to reference's
func SameVideoInfo(refPath, disPath string) (bool, VideoInfo, error) {
	// Get video information
	refInfo, err := GetVideoInfo(refPath)
	if err != nil {
		return false, VideoInfo{}, err
	}

	disInfo, err := GetVideoInfo(disPath)
	if err != nil {
		return false, VideoInfo{}, err
	}

	// Compare frame rate strings
	refFPS, _ := parseFPS(refInfo.FrameRate)
	disFPS, _ := parseFPS(disInfo.FrameRate)
	if refFPS != disFPS {
		return false, VideoInfo{}, fmt.Errorf(
			"reference and distorted files have different framerates: %s, %s",
			refInfo.FrameRate,
			disInfo.FrameRate,
		)
	}

	// Compare frame count
	if refInfo.FrameCount != disInfo.FrameCount {
		return false, VideoInfo{}, fmt.Errorf(
			"reference and distorted files have different frame counts: %s, %s",
			refInfo.FrameCount,
			disInfo.FrameCount,
		)
	}

	return true, refInfo, nil
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
		"-show_entries", "stream=width,height,r_frame_rate,nb_frames,pix_fmt",
		"-of", "json",
		path,
	)
	hideCmdWindow(cmd)

	// Get command output
	out, err := cmd.CombinedOutput()
	if err != nil {
		return VideoInfo{}, fmt.Errorf("%v: %s", err, string(out))
	}

	// Unmarshal to struct
	var res ffprobeOut
	if err := json.Unmarshal(out, &res); err != nil {
		return VideoInfo{}, err
	}

	// Ensure only one video stream exists
	if len(res.Streams) != 1 {
		return VideoInfo{}, fmt.Errorf("expected exactly one video stream, but found none or multiple in file: %s", path)
	}

	// Ensure resolution exists
	if res.Streams[0].Width == 0 || res.Streams[0].Height == 0 {
		return VideoInfo{}, fmt.Errorf("unable to get resolution of file: %s", path)
	}

	// Ensure frame rate is parse-able
	_, err = parseFPS(res.Streams[0].FrameRate)
	if err != nil {
		return VideoInfo{}, fmt.Errorf("%s of file: %s", err, path)
	}

	// Ensure frame count exists
	i, err := strconv.Atoi(res.Streams[0].FrameCount)
	if err != nil || i == 0 {
		// Fallback to counting frames directly
		cmd = exec.Command(
			"ffprobe",
			"-v", "error",
			"-select_streams", "v:0",
			"-count_frames",
			"-show_entries", "stream=nb_read_frames",
			"-of", "default=nokey=1:noprint_wrappers=1",
			path,
		)
		hideCmdWindow(cmd)

		out, err = cmd.CombinedOutput()
		if err != nil {
			return VideoInfo{}, fmt.Errorf("%v: %s", err, string(out))
		}

		i, err = strconv.Atoi(strings.TrimSpace(string(out)))
		if err != nil || i == 0 {
			return VideoInfo{}, fmt.Errorf("unable to get frame count of file: %s", path)
		}

		res.Streams[0].FrameCount = strconv.Itoa(i)
	}

	// Ensure pixel format exists
	if res.Streams[0].PixFmt == "" {
		return VideoInfo{}, fmt.Errorf("unable to get pixel format of file: %s", path)
	}

	return res.Streams[0], nil
}
