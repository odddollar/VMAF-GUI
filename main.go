package main

import (
	"VMAF-GUI/ui"
)

func main() {
	// ch, err := ffmpeg.RunVMAF("reference.mp4", "distorted.mp4")
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
