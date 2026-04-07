package ui

import (
	"VMAF-GUI/video"
	"context"

	"fyne.io/fyne/v2"
)

// Get frame from index and update compare widget
func (u *Ui) updateCompareImageIndex(index int) {
	// Get frames
	refImg, disImg, err := video.GetFramePair(
		context.TODO(),
		u.referenceEntry.Text,
		u.distortedEntry.Text,
		u.refInfo,
		index,
	)
	if err != nil {
		u.showError(err, false)
		return
	}

	fyne.Do(func() {
		// Update compare image
		u.compareImages.SetImages(refImg, disImg)
	})
}
