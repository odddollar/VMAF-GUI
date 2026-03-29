package ui

import (
	"VMAF-GUI/video"
	"fmt"
)

// Perform checks and run vmaf calculation
func (u *Ui) run() {
	// Ensure matching video info
	same, err := video.SameVideoInfo(u.referenceEntry.Text, u.distortedEntry.Text)
	if !same || err != nil {
		u.showError(err, false)
		return
	}

	fmt.Println("Running")
}
