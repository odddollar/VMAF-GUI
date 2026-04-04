package widgets

import (
	"VMAF-GUI/video"
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// Generate empty image
func newEmptyImage(width, height int, c color.Color) *image.NRGBA {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	// Fill background
	for y := range height {
		for x := range width {
			img.Set(x, y, c)
		}
	}

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
	widget *VMAFGraph
	raster *canvas.Raster
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

	return out
}
