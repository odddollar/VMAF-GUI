package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	// Main app elements
	a fyne.App
	w fyne.Window

	// Main UI elements
	titleLabel      *canvas.Text
	referenceEntry  *widget.Entry
	distortedEntry  *widget.Entry
	referenceButton *widget.Button
	distortedButton *widget.Button
	startButton     *widget.Button
}

func (u *Ui) NewUI() {
	// Create window
	u.a = app.New()
	u.w = u.a.NewWindow("VMAF GUI")

	// Create title widget
	u.titleLabel = canvas.NewText("VMAF GUI", color.Black)
	u.titleLabel.Alignment = fyne.TextAlignCenter
	u.titleLabel.TextStyle.Bold = true
	u.titleLabel.TextSize = 20

	// Create file path widgets
	u.referenceEntry = widget.NewEntry()
	u.referenceEntry.Validator = validateFileExists
	u.referenceEntry.OnChanged = func(s string) {
		u.validatePathEntries()
	}
	u.distortedEntry = widget.NewEntry()
	u.distortedEntry.Validator = validateFileExists
	u.distortedEntry.OnChanged = func(s string) {
		u.validatePathEntries()
	}

	// Create file explore buttons
	u.referenceButton = widget.NewButton("...", func() { u.selectFile(u.referenceEntry) })
	u.distortedButton = widget.NewButton("...", func() { u.selectFile(u.distortedEntry) })

	// Create start button
	u.startButton = widget.NewButton("Run", func() {})
	u.startButton.Importance = widget.HighImportance
	u.startButton.Disable()

	// Create window layout and set content
	u.w.SetContent(container.NewVBox(
		u.titleLabel,
		widget.NewForm(
			widget.NewFormItem("Reference file",
				container.NewBorder(nil, nil, nil,
					u.referenceButton,
					u.referenceEntry,
				),
			),
			widget.NewFormItem("Distorted file",
				container.NewBorder(nil, nil, nil,
					u.distortedButton,
					u.distortedEntry,
				),
			),
		),
		u.startButton,
	))
}

func (u *Ui) Run() {
	u.w.Resize(fyne.NewSize(800, 600))
	u.w.Show()
	u.a.Run()
}

// Disables start button if paths are invalid
func (u *Ui) validatePathEntries() {
	referenceErr := u.referenceEntry.Validate()
	distortedErr := u.distortedEntry.Validate()

	if referenceErr == nil && distortedErr == nil {
		fyne.Do(func() { u.startButton.Enable() })
	} else {
		fyne.Do(func() { u.startButton.Disable() })
	}
}
