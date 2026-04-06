package widgets

import (
	"VMAF-GUI/video"
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Generate empty image
func newEmptyImage(width, height int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Fill background
	draw.Draw(img, img.Bounds(), &image.Uniform{C: c}, image.Point{}, draw.Src)

	return img
}

// Get absolute value
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Draw line between two sets of coordinates
func drawLine(img *image.NRGBA, x0, y0, x1, y1 int, c color.Color) {
	// Calculate distance between points
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)

	// Step direction for x
	sx := -1
	if x0 < x1 {
		sx = 1
	}

	// Step direction for y
	sy := -1
	if y0 < y1 {
		sy = 1
	}

	// Total error/distance from destination
	err := dx + dy

	for {
		// Draw current pixel if inside bounds
		if image.Pt(x0, y0).In(img.Bounds()) {
			img.Set(x0, y0, c)
		}

		// Endpoint reached
		if x0 == x1 && y0 == y1 {
			break
		}

		// Double error for comparison
		e2 := 2 * err

		// Move in x direction if error threshold crossed
		if e2 >= dy {
			err += dy
			x0 += sx
		}

		// Move in y direction if error threshold crossed
		if e2 <= dx {
			err += dx
			y0 += sy
		}
	}
}

// Custom widget that displays vmaf results across video frames
type VMAFGraph struct {
	widget.BaseWidget

	// Hold vmaf data for graphing
	vmafData video.VMAFOutput

	// Update tooltip position
	mousePos fyne.Position
	mouseIn  bool
}

// Creates new VMAFGraph widget
func NewVMAFGraph() *VMAFGraph {
	// Create new object
	graph := &VMAFGraph{}

	graph.ExtendBaseWidget(graph)
	return graph
}

// Updates tooltip when mouse enters widget
func (w *VMAFGraph) MouseIn(event *desktop.MouseEvent) {
	w.mousePos = event.Position
	w.mouseIn = true
	w.Refresh()
}

// Updates tooltip when mouse moves over widget
func (w *VMAFGraph) MouseMoved(event *desktop.MouseEvent) {
	w.mousePos = event.Position
	w.Refresh()
}

// Hides tooltip when mouse leaves widget
func (w *VMAFGraph) MouseOut() {
	w.mouseIn = false
	w.Refresh()
}

// Set vmaf results
func (w *VMAFGraph) SetVMAF(vmaf video.VMAFOutput) {
	w.vmafData = vmaf
	w.Refresh()
}

// Returns new renderer for VMAFGraph
func (w *VMAFGraph) CreateRenderer() fyne.WidgetRenderer {
	r := &vmafGraphRenderer{}
	r.widget = w
	r.raster = canvas.NewRaster(r.generate)
	return r
}

// Renderer for VMAFGraph widget
type vmafGraphRenderer struct {
	widget   *VMAFGraph
	raster   *canvas.Raster
	fontFace font.Face
}

// Returns minimum size of VMAFGraph
func (r *vmafGraphRenderer) MinSize() fyne.Size {
	return fyne.NewSize(250, 250)
}

// Lays out raster to fill VMAFGraph
func (r *vmafGraphRenderer) Layout(size fyne.Size) {
	r.raster.Resize(size)
}

// Refreshes VMAFGraph
func (r *vmafGraphRenderer) Refresh() {
	r.raster.Refresh()
}

// Returns child widgets of VMAFGraph
func (r *vmafGraphRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.raster}
}

// Does nothing as VMAFGraph doesn't hold external resources
func (r *vmafGraphRenderer) Destroy() {}

// Generates composite RGBA image for raster callback
func (r *vmafGraphRenderer) generate(width, height int) image.Image {
	// Create new image to draw to
	out := newEmptyImage(width, height, color.Black)
	if width == 0 || height == 0 {
		return out
	}

	// Get frame data
	frames := r.widget.vmafData.Frames
	n := len(frames)
	if n == 0 {
		return out
	}

	// Graph area
	gw := float32(width)
	gh := float32(height)

	// Convert frame index to x coordinate
	scaleX := func(i int) int {
		return int((float32(i) / float32(n-1)) * gw)
	}

	// Convert vmaf score to y coordinate
	scaleY := func(v float64) int {
		return int((1.0 - float32(v)/100.0) * gh)
	}

	// Draw single line if only one frame
	if n == 1 {
		y := int((1.0 - float32(frames[0].Metrics.VMAF)/100.0) * gh)

		drawLine(out, 0, y, width-1, y, theme.Color(theme.ColorNamePrimary))
	} else {
		// Draw lines for every vmaf frame value
		for i := 0; i < n-1; i++ {
			x1 := scaleX(i)
			y1 := scaleY(frames[i].Metrics.VMAF)

			x2 := scaleX(i + 1)
			y2 := scaleY(frames[i+1].Metrics.VMAF)

			drawLine(out, x1, y1, x2, y2, theme.Color(theme.ColorNamePrimary))
		}
	}

	// Only draw tooltip if mouse in and frames exist
	if r.widget.mouseIn && n > 0 {
		// Initialise font for tooltip drawing
		r.initialiseFont()

		// Scale mouse position to frame number
		i := min(
			max(
				int((r.widget.mousePos.X/gw)*float32(n-1)+0.5),
				0,
			),
			n-1,
		)

		// Create text
		frame := frames[i]
		lines := []string{
			fmt.Sprintf("Frame: %d", frame.FrameNum+1),
			fmt.Sprintf("VMAF: %.2f", frame.Metrics.VMAF),
		}

		// Get proper tooltip sizing
		padding := float32(6)
		offset := float32(12)

		// Measure text using font metrics
		var textWidth float32
		var textHeight float32
		var lineHeight float32
		for _, line := range lines {
			bounds, _ := font.BoundString(r.fontFace, line)

			w := float32((bounds.Max.X - bounds.Min.X).Ceil())
			h := float32((bounds.Max.Y - bounds.Min.Y).Ceil())

			if w > textWidth {
				textWidth = w
			}
			lineHeight = h
		}
		textHeight = lineHeight * float32(len(lines))
		bgWidth := textWidth + padding*2 - 1
		bgHeight := textHeight + padding*float32(len(lines)+1) - 2

		// Put tooltip in bottom right corner of cursor
		tx := r.widget.mousePos.X + offset
		ty := r.widget.mousePos.Y + offset

		// Flip position horizontally
		if tx+bgWidth > float32(width) {
			tx = r.widget.mousePos.X - bgWidth
		}

		// Flip position vertically
		if ty+bgHeight > float32(height) {
			ty = r.widget.mousePos.Y - bgHeight
		}

		// Draw background to image
		bgRect := image.Rect(
			int(tx),
			int(ty),
			int(tx+bgWidth),
			int(ty+bgHeight),
		)
		draw.Draw(out, bgRect, &image.Uniform{
			C: color.RGBA{32, 32, 36, 235},
		}, image.Point{}, draw.Over)

		// Draw text to image over background
		d := &font.Drawer{
			Dst:  out,
			Src:  image.NewUniform(color.White),
			Face: r.fontFace,
		}
		textY := ty + padding + lineHeight
		for _, line := range lines {
			d.Dot = fixed.Point26_6{
				X: fixed.I(int(tx + padding - 1)),
				Y: fixed.I(int(textY - 1)),
			}
			d.DrawString(line)
			textY += lineHeight + padding
		}
	}

	return out
}

// Initialise font for text rendering
func (r *vmafGraphRenderer) initialiseFont() {
	// Don't initialise more than once
	if r.fontFace != nil {
		return
	}

	// Get fyne theme's font
	fontRes := theme.TextMonospaceFont()

	// Parse font
	tt, err := opentype.Parse(fontRes.Content())
	if err != nil {
		return
	}

	// Create usable font face
	face, err := opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    13,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return
	}

	r.fontFace = face
}
