package ui

import (
	"VMAF-GUI/video"
	"context"
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

	// Switch which button visible
	u.showStopButton()

	// Create context to allow vmaf command cancelling
	ctx, cancel := context.WithCancel(context.Background())
	u.vmafCancel = cancel

	// Start vmaf with channels
	progressChan, errChan, err := video.RunVMAF(ctx, u.referenceEntry.Text, u.distortedEntry.Text)
	if err != nil {
		u.showError(err, false)
		return
	}

	// Read channels
	go func() {
		for {
			select {
			case progress, ok := <-progressChan: // Get progress update
				if !ok {
					return
				}

				fmt.Println(progress)

			case err, ok := <-errChan: // Handle errors
				if !ok {
					return
				}

				u.showError(err, false)
				u.showStartButton()

				// Cancel vmaf calculation
				if u.vmafCancel != nil {
					u.vmafCancel()
					u.vmafCancel = nil
				}
				return
			}
		}
	}()
}

// Stop running vmaf calculation
func (u *Ui) stop() {
	u.vmafCancel()
	u.showStartButton()
}
