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
	return fyne.NewSize(200, 200)
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

	// Only draw tooltip if mouse in
	if r.widget.mouseIn {
		// Initialise font for tooltip drawing
		r.initialiseFont()

		// Create text
		label := fmt.Sprintf("%.0f, %.0f",
			r.widget.mousePos.X,
			r.widget.mousePos.Y,
		)

		// Get proper tooltip sizing
		padding := float32(6)
		offset := float32(12)

		// Measure text using font metrics
		bounds, _ := font.BoundString(r.fontFace, label)
		textWidth := float32((bounds.Max.X - bounds.Min.X).Ceil())
		textHeight := float32((bounds.Max.Y - bounds.Min.Y).Ceil())
		bgWidth := textWidth + padding*2
		bgHeight := textHeight + padding*2

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
			Dot: fixed.Point26_6{
				X: fixed.I(int(tx + padding - 1)),
				Y: fixed.I(int(ty + padding + textHeight - 1)),
			},
		}
		d.DrawString(label)
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
