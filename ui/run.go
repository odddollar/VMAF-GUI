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
	u.disableRunningWidgets()
	u.disableBottomWidgets()
	u.showStopButton()
	u.resetState()

	// Get reference info to update progress bar maximum
	u.refInfo, err = video.GetVideoInfo(u.referenceEntry.Text)
	if err != nil {
		u.showErrorAndReset(err, false)
		return
	}
	frameCount, err := strconv.ParseFloat(u.refInfo.FrameCount, 64)
	if err != nil {
		u.showErrorAndReset(err, false)
		return
	}
	u.progressBar.Max = frameCount
	u.maxFrameBinding.Set(int(frameCount))

	// Create context to allow vmaf command cancelling
	ctx, cancel := context.WithCancel(context.Background())
	u.vmafCancel = cancel

	// Intercept to stop calculation when app closed
	u.w.SetCloseIntercept(func() {
		if u.vmafCancel != nil {
			u.vmafCancel()
			u.vmafCancel = nil
		}
		u.w.Close()
	})

	// Start vmaf with channels
	progressChan, errChan, doneChan, err := video.RunVMAF(ctx, u.referenceEntry.Text, u.distortedEntry.Text)
	if err != nil {
		u.showErrorAndReset(err, false)
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
				u.showErrorAndReset(err, false)

				// Cancel vmaf calculation
				if u.vmafCancel != nil {
					u.vmafCancel()
					u.vmafCancel = nil
				}
				return

			case <-doneChan: // Command finished successfully
				u.enableRunningWidgets()
				u.enableBottomWidgets()
				u.showStartButton()

				// Parse vmaf results and store
				vmaf, err := video.ParseJsonOutput("vmaf.json", u.deleteOutputCheck.Checked)
				if err != nil {
					u.showError(err, false)
					return
				}

				// Update compare widget to first frame
				go u.updateCompareImageIndex(0)

				fyne.Do(func() {
					// Update results
					u.resultsMeanBinding.Set(vmaf.PooledMetrics.VMAF.Mean)
					u.resultsHarmonicMeanBinding.Set(vmaf.PooledMetrics.VMAF.HarmonicMean)
					u.resultsMinBinding.Set(vmaf.PooledMetrics.VMAF.Min)
					u.resultsMaxBinding.Set(vmaf.PooledMetrics.VMAF.Max)

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
	if u.vmafCancel != nil {
		u.vmafCancel()
		u.vmafCancel = nil
	}
	u.enableRunningWidgets()
	u.showStartButton()
	u.resetState()
}
