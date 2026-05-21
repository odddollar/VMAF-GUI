package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// Use fyne.Do() for errors as they can occur in any thread

// Standard dialog to show error
func (u *Ui) showError(err error, fatal bool) {
	fyne.Do(func() {
		d := dialog.NewError(err, u.w)

		if fatal {
			// Close window if error fatal
			d.SetOnClosed(func() {
				u.a.Quit()
			})

			// Log fatal error
			u.logger.Printf("FATAL: %v", err)
		} else {
			// Log regular error
			u.logger.Printf("ERROR: %v", err)
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
