package ui

import (
	"VMAF-GUI/widgets"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/theme"
)

// Create vmaf scale with given text alignment
func newVmafScale(alignment fyne.TextAlign) *fyne.Container {
	top := canvas.NewText("100", theme.Color(theme.ColorNameForeground))
	top.Alignment = alignment
	top.TextStyle.Bold = true

	mid := canvas.NewText("VMAF", theme.Color(theme.ColorNameForeground))
	mid.Alignment = alignment
	mid.TextStyle.Bold = true

	bottom := canvas.NewText("0", theme.Color(theme.ColorNameForeground))
	bottom.Alignment = alignment
	bottom.TextStyle.Bold = true

	return container.NewBorder(
		top,
		bottom,
		nil,
		nil,
		mid,
	)
}

// Create frame scale with given max frame
func newFrameScale(maxFrameBinding binding.Int) *fyne.Container {
	// Alignment spacer
	t := canvas.NewText("VMAF", theme.Color(theme.ColorNameForeground))
	t.TextStyle.Bold = true
	spacer := widgets.NewSpacer(t.MinSize())

	left := canvas.NewText("1", theme.Color(theme.ColorNameForeground))
	left.Alignment = fyne.TextAlignLeading
	left.TextStyle.Bold = true

	middle := canvas.NewText("Frame number", theme.Color(theme.ColorNameForeground))
	middle.Alignment = fyne.TextAlignCenter
	middle.TextStyle.Bold = true

	right := canvas.NewText("0", theme.Color(theme.ColorNameForeground))
	right.Alignment = fyne.TextAlignTrailing
	right.TextStyle.Bold = true

	// Attach listener to binding
	maxFrameBinding.AddListener(binding.NewDataListener(func() {
		val, _ := maxFrameBinding.Get()
		right.Text = fmt.Sprintf("%d", val)
		right.Refresh()
	}))

	return container.NewGridWithColumns(3,
		container.NewBorder(nil, nil, spacer, nil, left),
		middle,
		container.NewBorder(nil, nil, nil, spacer, right),
	)
}
