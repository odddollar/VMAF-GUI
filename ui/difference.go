package ui

import (
	"VMAF-GUI/video"
	"context"
	"strconv"
	"unicode"

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

// Ensures only numbers less than max frame count entered
func (u *Ui) restrictCompareFrameEntry(s string) {
	filtered := ""

	// Restrict to digits only
	for _, r := range s {
		if unicode.IsDigit(r) {
			filtered += string(r)
		}
	}

	// If nothing valid default to 1
	if filtered == "" {
		filtered = "1"
	}

	val, err := strconv.Atoi(filtered)
	if err != nil {
		return
	}

	maxFrame, err := u.maxFrameBinding.Get()
	if err != nil {
		return
	}

	// Clamp between 1 and max
	val = min(max(val, 1), maxFrame)

	new := strconv.Itoa(val)

	// Update content if changed
	if s != new {
		u.compareFrameEntry.SetText(new)
	}
}
