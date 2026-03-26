package ui

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// Show dialog for selecting file
func (u *Ui) selectFile(target *widget.Entry) {
	d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			u.showError(err, false)
			return
		}

		// If nothing selected
		if reader == nil {
			return
		}

		// Close reader
		defer reader.Close()

		// Update entry path
		target.SetText(reader.URI().Path())
	}, u.w)

	// Filter for video file types
	d.SetFilter(storage.NewExtensionFileFilter([]string{
		".mp4", ".mkv", ".avi", ".mov", ".wmv",
		".flv", ".webm", ".mpeg", ".mpg", ".m4v",
		".3gp", ".ogv", ".ts", ".mts", ".vob",
	}))

	// Open to current directory
	// Should fallback to user directory
	cwd, err := os.Getwd()
	if err == nil {
		// Convert to uri
		uri := storage.NewFileURI(cwd)
		listable, err := storage.ListerForURI(uri)
		if err == nil {
			d.SetLocation(listable)
		}
	}

	d.SetTitleText("Select video file")
	d.SetConfirmText("Select")
	d.SetView(dialog.ListView)
	d.Show()
}
