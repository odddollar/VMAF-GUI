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

		// Close window if error fatal
		if fatal {
			d.SetOnClosed(func() {
				u.a.Quit()
			})
		}

		d.Show()
	})
}
