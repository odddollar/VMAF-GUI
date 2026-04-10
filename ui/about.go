package ui

import (
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/dialog"
)

// Use Fyne-X extensions to create about window
func (u *Ui) showAbout() {
	// Parse urls
	github, _ := url.Parse("https://github.com/odddollar/VMAF-GUI")

	links := []*widget.Hyperlink{
		widget.NewHyperlink("VMAF GUI GitHub", github),
	}

	// Markdown program description
	content := "A UI for Netflix's **VMAF** video comparison algorithm"

	// Use Fyne-X's about dialog
	d := dialog.NewAbout(content, links, u.a, u.w)
	d.Resize(fyne.NewSize(400, 406))
	d.Show()
}
