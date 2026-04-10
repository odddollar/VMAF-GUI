package ui

import (
	"VMAF-GUI/video"
	"VMAF-GUI/widgets"
	"context"
	"fmt"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	// Main app elements
	a fyne.App
	w fyne.Window

	// Global bindings
	maxFrameBinding binding.Int

	// Bindings for progress tracking
	progressFrameBinding   binding.Int
	progressFpsBinding     binding.Int
	progressElapsedBinding binding.String

	// Bindings for results display
	resultsMeanBinding         binding.Float
	resultsHarmonicMeanBinding binding.Float
	resultsMinBinding          binding.Float
	resultsMaxBinding          binding.Float

	// Bindings for compare vmaf display
	compareVmafBinding binding.Float

	// Main UI elements
	titleLabel        *canvas.Text
	aboutButton       *widget.Button
	referenceEntry    *widget.Entry
	distortedEntry    *widget.Entry
	referenceButton   *widget.Button
	distortedButton   *widget.Button
	modelDropdown     *widget.Select
	modelInfoButton   *widget.Button
	startButton       *widget.Button
	stopButton        *widget.Button
	deleteOutputCheck *widget.Check
	progressBar       *widget.ProgressBar
	progressLabel     *widget.Label

	// Results tab elements
	resultsLabel            *widget.Label
	resultsGraph            *widgets.VMAFGraph
	resultsLeftVMAFLabels   *fyne.Container
	resultsRightVMAFLabels  *fyne.Container
	resultsFrameCountLabels *fyne.Container

	// Compare tab elements
	compareResultsLabel *widget.Label
	compareImages       *widgets.CompareWidget
	comparePrevButton   *widget.Button
	compareNextButton   *widget.Button
	compareFrameEntry   *widget.Entry
	compareFrameLabel   *widget.Label

	// Allows cancelling in-progress vmaf calculation
	vmafCancel context.CancelFunc

	// Allows cancelling in-progress frame extraction
	compareCancel    context.CancelFunc
	compareRequestId int

	// Store information for current reference file
	refInfo video.VideoInfo

	// Store most recent vmaf scores
	vmafScores video.VMAFOutput
}

func (u *Ui) NewUI() {
	// Create window
	u.a = app.New()
	u.w = u.a.NewWindow("VMAF GUI")

	// Create bindings for frame number tracking
	u.maxFrameBinding = binding.NewInt()
	u.maxFrameBinding.Set(0)

	// Create bindings for progress status
	u.progressFrameBinding = binding.NewInt()
	u.progressFpsBinding = binding.NewInt()
	u.progressElapsedBinding = binding.NewString()
	u.progressElapsedBinding.Set("0s")

	// Create bindings for results
	u.resultsMeanBinding = binding.NewFloat()
	u.resultsHarmonicMeanBinding = binding.NewFloat()
	u.resultsMinBinding = binding.NewFloat()
	u.resultsMaxBinding = binding.NewFloat()

	// Create bindings for compare vmaf
	u.compareVmafBinding = binding.NewFloat()

	// Create title widget
	u.titleLabel = canvas.NewText("VMAF GUI", theme.Color(theme.ColorNameForeground))
	u.titleLabel.Alignment = fyne.TextAlignCenter
	u.titleLabel.TextStyle.Bold = true
	u.titleLabel.TextSize = 20

	// Create about button
	u.aboutButton = widget.NewButtonWithIcon("", theme.InfoIcon(), u.showAbout)

	// Create file path widgets
	u.referenceEntry = widget.NewEntry()
	u.referenceEntry.Validator = validateFileExists
	u.referenceEntry.OnChanged = func(s string) {
		u.validatePathEntries()
		u.disableBottomWidgets()
	}
	u.distortedEntry = widget.NewEntry()
	u.distortedEntry.Validator = validateFileExists
	u.distortedEntry.OnChanged = func(s string) {
		u.validatePathEntries()
		u.disableBottomWidgets()
	}

	// Create file explore buttons
	u.referenceButton = widget.NewButtonWithIcon("Browse", theme.SearchIcon(), func() { u.selectFile(u.referenceEntry) })
	u.distortedButton = widget.NewButtonWithIcon("Browse", theme.SearchIcon(), func() { u.selectFile(u.distortedEntry) })

	// Create model selection dropdown
	u.modelDropdown = widget.NewSelect([]string{
		"vmaf_v0.6.1",
		"vmaf_4k_v0.6.1",
		"vmaf_v0.6.1neg",
		"vmaf_4k_v0.6.1neg",
	}, func(s string) {})
	u.modelDropdown.SetSelectedIndex(0)

	// Create model info button
	u.modelInfoButton = widget.NewButtonWithIcon("Models", theme.InfoIcon(), u.showModelInfo)

	// Create start button
	u.startButton = widget.NewButton("Start", func() { go u.run() })
	u.startButton.Importance = widget.HighImportance
	u.startButton.Disable()

	// Create stop button
	u.stopButton = widget.NewButton("Stop", func() { go u.stop() })
	u.stopButton.Hide()

	// Create delete output check box
	u.deleteOutputCheck = widget.NewCheck("Delete \"vmaf.json\"", func(b bool) {})
	u.deleteOutputCheck.SetChecked(true)

	// Create progress bar
	u.progressBar = widget.NewProgressBar()

	// Create progress label
	progressStatus := binding.NewSprintf(
		"Frame: %4d, FPS: %3d, Elapsed: %6s",
		u.progressFrameBinding,
		u.progressFpsBinding,
		u.progressElapsedBinding,
	)
	u.progressLabel = widget.NewLabelWithData(progressStatus)
	u.progressLabel.TextStyle.Monospace = true

	// Top main UI elements
	topElements := container.NewVBox(
		container.NewBorder(
			nil,
			nil,
			widgets.NewSpacer(widget.NewButtonWithIcon("", theme.InfoIcon(), func() {}).MinSize()), // Keeps title centred
			u.aboutButton,
			u.titleLabel,
		),
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
			widget.NewFormItem("VMAF model",
				container.NewBorder(nil, nil, nil,
					u.modelInfoButton,
					u.modelDropdown,
				),
			),
		),
		container.NewBorder(nil, nil, nil,
			u.deleteOutputCheck,
			container.NewVBox(
				u.startButton,
				u.stopButton,
			),
		),
		container.NewBorder(nil, nil, nil,
			u.progressLabel,
			u.progressBar,
		),
	)

	// Create results label
	resultsText := binding.NewSprintf(
		"VMAF mean: %.2f, VMAF harmonic mean: %.2f, VMAF minimum: %.2f, VMAF maximum: %.2f",
		u.resultsMeanBinding,
		u.resultsHarmonicMeanBinding,
		u.resultsMinBinding,
		u.resultsMaxBinding,
	)
	u.resultsLabel = widget.NewLabelWithData(resultsText)
	u.resultsLabel.TextStyle.Monospace = true

	// Create results graph
	u.resultsGraph = widgets.NewVMAFGraph()

	// Create results graph labels
	u.resultsLeftVMAFLabels = newVmafScale(fyne.TextAlignTrailing)
	u.resultsRightVMAFLabels = newVmafScale(fyne.TextAlignLeading)
	u.resultsFrameCountLabels = newFrameScale(u.maxFrameBinding)

	// Results tab elements
	resultsTabElements := container.NewBorder(
		container.NewHBox(
			layout.NewSpacer(),
			u.resultsLabel,
			layout.NewSpacer(),
		),
		u.resultsFrameCountLabels,
		u.resultsLeftVMAFLabels,
		u.resultsRightVMAFLabels,
		u.resultsGraph,
	)

	// Create compare vmaf label
	compareVmafText := binding.NewSprintf(
		"Frame VMAF: %.2f",
		u.compareVmafBinding,
	)
	u.compareResultsLabel = widget.NewLabelWithData(compareVmafText)
	u.compareResultsLabel.TextStyle.Monospace = true

	// Create compare image widget
	u.compareImages = widgets.NewCompareWidget(image.Black, image.Black)

	// Create previous and next frame buttons
	u.comparePrevButton = widget.NewButtonWithIcon("", theme.NavigateBackIcon(), u.compareFrameEntryPrev)
	u.compareNextButton = widget.NewButtonWithIcon("", theme.NavigateNextIcon(), u.compareFrameEntryNext)

	// Create entry that only allows digits
	u.compareFrameEntry = widget.NewEntry()
	u.compareFrameEntry.SetText("1")
	u.compareFrameEntry.OnChanged = u.compareFrameEntryRestrict

	// Create dynamic label for max frame number
	u.compareFrameLabel = widget.NewLabelWithData(binding.NewSprintf(
		"of %d",
		u.maxFrameBinding,
	))
	u.compareFrameLabel.TextStyle.Bold = true

	// Compare tab elements
	compareTabElements := container.NewBorder(
		container.NewHBox(
			layout.NewSpacer(),
			u.compareResultsLabel,
			layout.NewSpacer(),
		),
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
		u.compareImages,
	)

	// Create window layout and set content
	u.w.SetContent(container.NewBorder(
		topElements,
		nil, nil, nil,
		container.NewAppTabs(
			container.NewTabItemWithIcon("Results", theme.ListIcon(), resultsTabElements),
			container.NewTabItemWithIcon("Compare", theme.VisibilityIcon(), compareTabElements),
		),
	))

	u.disableBottomWidgets()
}

func (u *Ui) Run() {
	u.w.Resize(fyne.NewSize(800, 0))
	u.w.Show()
	go u.startupChecks()
	u.a.Run()
}

// Checks to ensure program can run properly
func (u *Ui) startupChecks() {
	if !video.CommandAvailable("ffmpeg") {
		u.showError(fmt.Errorf("unable to find FFmpeg"), true)
		return
	}

	if !video.VMAFAvailable() {
		u.showError(fmt.Errorf("unable to find VMAF in FFmpeg"), true)
		return
	}

	if !video.CommandAvailable("ffprobe") {
		u.showError(fmt.Errorf("unable to find FFprobe"), true)
		return
	}
}

// Disable widgets that shouldn't be changed when running
func (u *Ui) disableRunningWidgets() {
	fyne.Do(func() {
		u.referenceEntry.Disable()
		u.distortedEntry.Disable()
		u.referenceButton.Disable()
		u.distortedButton.Disable()
		u.modelDropdown.Disable()
		u.modelInfoButton.Disable()
		u.deleteOutputCheck.Disable()
	})
}

// Enable widgets that shouldn't be changed when running
func (u *Ui) enableRunningWidgets() {
	fyne.Do(func() {
		u.referenceEntry.Enable()
		u.distortedEntry.Enable()
		u.referenceButton.Enable()
		u.distortedButton.Enable()
		u.modelDropdown.Enable()
		u.modelInfoButton.Enable()
		u.deleteOutputCheck.Enable()
	})
}

// Disable bottom widgets that shouldn't be enabled until successfully completed
func (u *Ui) disableBottomWidgets() {
	fyne.Do(func() {
		u.comparePrevButton.Disable()
		u.compareNextButton.Disable()
		u.compareFrameEntry.Disable()
	})
}

// Enable bottom widgets that shouldn't be enabled until successfully completed
func (u *Ui) enableBottomWidgets() {
	fyne.Do(func() {
		u.comparePrevButton.Enable()
		u.compareNextButton.Enable()
		u.compareFrameEntry.Enable()
	})
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

// Reset state of widgets to new
func (u *Ui) resetState() {
	fyne.Do(func() {
		u.progressBar.SetValue(0)
		u.progressFrameBinding.Set(0)
		u.progressFpsBinding.Set(0)
		u.progressElapsedBinding.Set("0s")
		u.compareFrameEntry.SetText("1")
	})
}
