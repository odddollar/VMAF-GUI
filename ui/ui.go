package ui

import (
	"VMAF-GUI/video"
	"errors"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
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
	progressBar     *widget.ProgressBar

	// Results tab elements

	// Compare tab elements
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
	u.referenceButton = widget.NewButtonWithIcon("Browse", theme.SearchIcon(), func() { u.selectFile(u.referenceEntry) })
	u.distortedButton = widget.NewButtonWithIcon("Browse", theme.SearchIcon(), func() { u.selectFile(u.distortedEntry) })

	// Create start button
	u.startButton = widget.NewButton("Run", func() {})
	u.startButton.Importance = widget.HighImportance
	u.startButton.Disable()

	// Create progress bar
	u.progressBar = widget.NewProgressBar()

	// Top main UI elements
	topElements := container.NewVBox(
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
		u.progressBar,
	)

	// Results tab elements
	resultsTabElements := container.NewVBox()

	// Compare tab elements
	compareTabElements := container.NewVBox()

	// Create window layout and set content
	u.w.SetContent(container.NewVBox(
		topElements,
		container.NewAppTabs(
			container.NewTabItemWithIcon("Results", theme.ListIcon(), resultsTabElements),
			container.NewTabItemWithIcon("Compare", theme.VisibilityIcon(), compareTabElements),
		),
	))
}

func (u *Ui) Run() {
	u.w.Resize(fyne.NewSize(800, 600))
	u.w.Show()
	u.startupChecks()
	u.a.Run()
}

// Runs checks to ensure program can run properly
func (u *Ui) startupChecks() {
	if !video.CommandAvailable("ffmpeg") {
		u.showError(errors.New("unable to find FFmpeg"), true)
		return
	}

	if !video.VMAFAvailable() {
		u.showError(errors.New("unable to find VMAF in FFmpeg"), true)
		return
	}

	if !video.CommandAvailable("ffprobe") {
		u.showError(errors.New("unable to find FFprobe"), true)
		return
	}
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
