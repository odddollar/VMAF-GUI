package ui

import (
	"VMAF-GUI/video"
	"context"
	"strconv"

	"fyne.io/fyne/v2"
)

// Perform checks and run vmaf calculation
func (u *Ui) run() {
	// Ensure matching video info
	same, err := video.SameVideoInfo(u.referenceEntry.Text, u.distortedEntry.Text)
	if !same || err != nil {
		u.showError(err, false)
		return
	}

	// Switch which button visible and clear progress
	u.showStopButton()
	u.clearProgressStatus()

	// Get reference info to update progress bar maximum
	refInfo, err := video.GetVideoInfo(u.referenceEntry.Text)
	if err != nil {
		u.showError(err, false)
		return
	}
	frameCount, err := strconv.ParseFloat(refInfo.FrameCount, 64)
	if err != nil {
		u.showError(err, false)
		return
	}
	u.progressBar.Max = frameCount
	u.maxFrameBinding.Set(int(frameCount))

	// Create context to allow vmaf command cancelling
	ctx, cancel := context.WithCancel(context.Background())
	u.vmafCancel = cancel

	// Start vmaf with channels
	progressChan, errChan, doneChan, err := video.RunVMAF(ctx, u.referenceEntry.Text, u.distortedEntry.Text)
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

				fyne.Do(func() {
					// Update progress bar
					u.progressBar.SetValue(float64(progress.Frame))
				})

				// Update progress label
				u.progressFrameBinding.Set(progress.Frame)
				u.progressFpsBinding.Set(progress.FPS)
				u.progressElapsedBinding.Set(progress.Elapsed.String())

			case err := <-errChan: // Handle errors
				u.showError(err, false)
				u.showStartButton()
				u.clearProgressStatus()

				// Cancel vmaf calculation
				if u.vmafCancel != nil {
					u.vmafCancel()
					u.vmafCancel = nil
				}
				return

			case <-doneChan: // Command finished successfully
				u.enableBottomWidgets()
				u.showStartButton()

				// Parse vmaf results and store
				vmaf, err := video.ParseJsonOutput("vmaf.json", u.deleteOutputCheck.Checked)
				if err != nil {
					u.showError(err, false)
					return
				}

				fyne.Do(func() {
					// Update graph
					u.resultsGraph.SetVMAF(vmaf)
				})

				return
			}
		}
	}()
}

// Stop running vmaf calculation
func (u *Ui) stop() {
	u.vmafCancel()
	u.showStartButton()
	u.clearProgressStatus()
}
