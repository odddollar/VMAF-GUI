package ui

import (
	"VMAF-GUI/video"
	"context"
	"strconv"
	"unicode"

	"fyne.io/fyne/v2"
)

// Get frame from index and update compare widget
func (u *Ui) compareImageUpdateIndex(index int) {
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
func (u *Ui) compareFrameEntryRestrict(s string) {
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

	// Clamp between 1 and max
	val, _ := strconv.Atoi(filtered)
	maxFrame, _ := u.maxFrameBinding.Get()
	val = min(max(val, 1), maxFrame)

	new := strconv.Itoa(val)

	// Update content if changed
	if s != new {
		u.compareFrameEntry.SetText(new)
	}
}

// Navigate to next frame
func (u *Ui) compareFrameEntryNext() {
	val := u.compareFrameEntry.Text
	valInt, _ := strconv.Atoi(val)
	valInt++
	val = strconv.Itoa(valInt)
	u.compareFrameEntry.SetText(val)
}

// Navigate to previous frame
func (u *Ui) compareFrameEntryPrev() {
	val := u.compareFrameEntry.Text
	valInt, _ := strconv.Atoi(val)
	valInt--
	val = strconv.Itoa(valInt)
	u.compareFrameEntry.SetText(val)
}
