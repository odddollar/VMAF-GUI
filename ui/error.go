package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

// Standard dialog to show error
func (u *Ui) showError(err error) {
	fyne.Do(func() { dialog.ShowError(err, u.w) })
}
