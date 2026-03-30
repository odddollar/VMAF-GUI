package ui

import (
	"VMAF-GUI/video"
	"VMAF-GUI/widgets"
	"context"
	"errors"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
	stopButton      *widget.Button
	progressBar     *widget.ProgressBar
	progressLabel   *widget.Label

	// Results tab elements

	// Compare tab elements
	compareImage      *widgets.CompareWidget
	comparePrevButton *widget.Button
	compareNextButton *widget.Button
	compareFrameEntry *widget.Entry
	compareFrameLabel *widget.Label

	// Bindings for progress tracking
	frameBinding   binding.Int
	fpsBinding     binding.Int
	elapsedBinding binding.String

	// Allows cancelling in-progress vmaf calculation
	vmafCancel context.CancelFunc
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
	u.startButton = widget.NewButton("Start", func() { go u.run() })
	u.startButton.Importance = widget.HighImportance
	u.startButton.Disable()

	// Create stop button
	u.stopButton = widget.NewButton("Stop", func() { go u.stop() })
	u.stopButton.Hide()

	// Create progress bar
	u.progressBar = widget.NewProgressBar()

	// Create bindings for progress status
	u.frameBinding = binding.NewInt()
	u.fpsBinding = binding.NewInt()
	u.elapsedBinding = binding.NewString()
	u.elapsedBinding.Set("0s")

	// Create progress label
	progressStatus := binding.NewSprintf(
		"Frame: %4d, FPS: %3d, Elapsed: %6s",
		u.frameBinding,
		u.fpsBinding,
		u.elapsedBinding,
	)
	u.progressLabel = widget.NewLabelWithData(progressStatus)
	u.progressLabel.TextStyle.Monospace = true

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
		u.stopButton,
		container.NewBorder(nil, nil, nil,
			u.progressLabel,
			u.progressBar,
		),
	)

	// Results tab elements
	resultsTabElements := container.NewVBox()

	// Create compare image widget
	u.compareImage = widgets.NewCompareWidget(image.Black, image.NewUniform(color.Gray{120}))

	// Create previous and next frame buttons
	u.comparePrevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {})
	u.compareNextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {})

	// Create entry and label for frame number tracking
	u.compareFrameEntry = widget.NewEntry()
	u.compareFrameLabel = widget.NewLabel("of 4096")
	u.compareFrameLabel.TextStyle.Bold = true

	// Compare tab elements
	compareTabElements := container.NewBorder(
		nil,
		container.NewCenter(container.NewHBox(
			u.comparePrevButton,
			container.NewGridWrap( // Expand width of entry
				fyne.NewSize(65, u.compareFrameEntry.MinSize().Height),
				u.compareFrameEntry,
			),
			u.compareFrameLabel,
			u.compareNextButton,
		)),
		nil,
		nil,
		u.compareImage,
	)

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
	u.w.Resize(fyne.NewSize(800, 0))
	u.w.Show()
	u.startupChecks()
	u.a.Run()
}

// Checks to ensure program can run properly
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

// Show start button
func (u *Ui) showStartButton() {
	fyne.Do(func() {
		u.startButton.Show()
		u.stopButton.Hide()
	})
}

// Show stop button
func (u *Ui) showStopButton() {
	fyne.Do(func() {
		u.startButton.Hide()
		u.stopButton.Show()
	})
}

// Clear progress status
func (u *Ui) clearProgressStatus() {
	fyne.Do(func() {
		u.progressBar.SetValue(0)
		u.frameBinding.Set(0)
		u.fpsBinding.Set(0)
		u.elapsedBinding.Set("0s")
	})
}
