package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// Use fyne.Do() for errors as they can occur in any thread

// Standard dialog to show error
func (u *Ui) showError(err error, fatal bool) {
	// Log error skipping stack frames
	msg := fmt.Sprintf("ERROR: %v", err)
	if fatal {
		msg = fmt.Sprintf("FATAL: %v", err)
	}
	u.logger.Output(3, msg)

	fyne.Do(func() {
		d := dialog.NewError(err, u.w)

		if fatal {
			// Close window if error fatal
			d.SetOnClosed(func() {
				u.a.Quit()
			})
		}

		d.Show()
	})
}

// Show error and reset ui
func (u *Ui) showErrorAndReset(err error, fatal bool) {
	u.enableRunningWidgets()
	u.showStartButton()
	u.resetState()
	u.showError(err, fatal)
}
