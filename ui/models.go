package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// Show dialog with information about vmaf models
func (u *Ui) showModelInfo() {
	// Create markdown content
	// Separate widgets for better spacing
	content1 := widget.NewRichTextFromMarkdown(`
**Standard - vmaf_v0.6.1**

- General-purpose quality model
- Suitable for most SDR video comparisons
- Balanced and widely used baseline
- Recommended default for typical workflows
`)

	content2 := widget.NewRichTextFromMarkdown(`
**4K UHD - vmaf_4k_v0.6.1**

- Trained specifically for high-resolution (2160p) content
- Better reflects perceived quality on large displays
- Use when working primarily with UHD sources
- Not ideal for lower resolutions
`)

	content3 := widget.NewRichTextFromMarkdown(`
**No Enhancement Gain - vmaf_v0.6.1neg**

- Penalises artificial sharpening and enhancement
- Prevents inflated scores from post-processing tricks
- Useful for codec testing and tuning
- Scores may appear lower than the standard model
`)

	content4 := widget.NewRichTextFromMarkdown(`
**4K UHD & No Enhancement Gain - vmaf_4k_v0.6.1neg**

- Combines UHD optimisation with enhancement penalties
- Best for high-resolution encoding tests
- More strict and realistic for modern pipelines
- Slightly higher computational cost
`)

	content5 := widget.NewRichTextFromMarkdown(`
**General Guidance**

- Use *Standard* unless you have a specific reason not to
- Use *NEG* variants when evaluating encoder performance
- Use *4K* variants only for UHD content
`)
	content1.Wrapping = fyne.TextWrapBreak
	content2.Wrapping = fyne.TextWrapBreak
	content3.Wrapping = fyne.TextWrapBreak
	content4.Wrapping = fyne.TextWrapBreak
	content5.Wrapping = fyne.TextWrapBreak

	// Show content in scrolling dialog
	d := dialog.NewCustom(
		"VMAF Model Information",
		"OK",
		container.NewVScroll(container.NewVBox(
			content1,
			content2,
			content3,
			content4,
			content5,
		)),
		u.w,
	)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}
