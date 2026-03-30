package widgets

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Clamp value between 0 and 1
func clamp01(v float32) float32 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}

// Custom widget that displays two images split by draggable vertical bar
type CompareWidget struct {
	widget.BaseWidget

	// Source images
	imgLeft  image.Image
	imgRight image.Image

	// Divider position from 0 to 1
	splitPos float32
}

// Creates new CompareWidget with left and right images
func NewCompareWidget(imgLeft, imgRight image.Image) *CompareWidget {
	w := &CompareWidget{
		imgLeft:  imgLeft,
		imgRight: imgRight,
		splitPos: 0.5,
	}
	w.ExtendBaseWidget(w)
	return w
}

// Slides divider as pointer dragged
func (w *CompareWidget) Dragged(e *fyne.DragEvent) {
	sz := w.Size()
	if sz.Width == 0 {
		return
	}
	w.splitPos = clamp01(e.Position.X / sz.Width)
	w.Refresh()
}

// Required by fyne.Draggable
func (w *CompareWidget) DragEnd() {}

// Moves divider to tapped position
func (w *CompareWidget) Tapped(e *fyne.PointEvent) {
	sz := w.Size()
	if sz.Width == 0 {
		return
	}
	w.splitPos = clamp01(e.Position.X / sz.Width)
	w.Refresh()
}

// Updates both images
func (w *CompareWidget) SetImages(imgA, imgB image.Image) {
	w.imgLeft = imgA
	w.imgRight = imgB
	w.Refresh()
}

// Sets divider position between 0 and 1
func (w *CompareWidget) SetSplitPosition(pos float32) {
	w.splitPos = clamp01(pos)
	w.Refresh()
}

// Returns new renderer for CompareWidget
func (w *CompareWidget) CreateRenderer() fyne.WidgetRenderer {
	r := &compareWidgetRenderer{}
	r.widget = w
	r.raster = canvas.NewRaster(r.generate)
	return r
}

// Renderer for CompareWidget
type compareWidgetRenderer struct {
	widget *CompareWidget
	raster *canvas.Raster
}

// Returns minimum size of CompareWidget
func (r *compareWidgetRenderer) MinSize() fyne.Size {
	return fyne.NewSize(200, 200)
}

// Lays out raster to fill CompareWidget
func (r *compareWidgetRenderer) Layout(size fyne.Size) {
	r.raster.Resize(size)
}

// Refreshes CompareWidget
func (r *compareWidgetRenderer) Refresh() {
	r.raster.Refresh()
}

// Returns child widgets of CompareWidget
func (r *compareWidgetRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.raster}
}

// Does nothing as CompareWidget doesn't hold external resources
func (r *compareWidgetRenderer) Destroy() {}

// Generates composite RGBA image for raster callback
func (r *compareWidgetRenderer) generate(width, height int) image.Image {
	// Create new image to draw to
	out := image.NewRGBA(image.Rect(0, 0, width, height))
	if width == 0 || height == 0 {
		return out
	}

	// Get x position of divider
	splitX := int(float32(width) * r.widget.splitPos)

	// Nearest-neighbour sample from img scaled to width, height
	sample := func(img image.Image, px, py int) color.Color {
		if img == nil {
			return color.Transparent
		}
		b := img.Bounds()
		imgW, imgH := b.Dx(), b.Dy()
		if imgW == 0 || imgH == 0 {
			return color.Transparent
		}

		// Largest uniform scale that fits image inside canvas
		scaleW := float64(width) / float64(imgW)
		scaleH := float64(height) / float64(imgH)
		scale := scaleW
		if scaleH < scale {
			scale = scaleH
		}

		scaledW := int(float64(imgW) * scale)
		scaledH := int(float64(imgH) * scale)
		offX := (width - scaledW) / 2
		offY := (height - scaledH) / 2

		// Fill outside image
		if px < offX || px >= offX+scaledW || py < offY || py >= offY+scaledH {
			return theme.Color(theme.ColorNameBackground)
		}

		// Map back to source pixel
		sx := b.Min.X + (px-offX)*imgW/scaledW
		sy := b.Min.Y + (py-offY)*imgH/scaledH
		if sx >= b.Max.X {
			sx = b.Max.X - 1
		}
		if sy >= b.Max.Y {
			sy = b.Max.Y - 1
		}
		return img.At(sx, sy)
	}

	// Draw left and right image halves either side of splitX
	for py := 0; py < height; py++ {
		for px := 0; px < splitX; px++ {
			out.Set(px, py, sample(r.widget.imgLeft, px, py))
		}
		for px := splitX; px < width; px++ {
			out.Set(px, py, sample(r.widget.imgRight, px, py))
		}
	}

	// Draw vertical divider
	for x := splitX - 1; x <= splitX+1; x++ {
		if x >= 0 && x < width {
			for y := 0; y < height; y++ {
				out.Set(x, y, theme.Color(theme.ColorNamePrimary))
			}
		}
	}

	// Set pixel if within bounds
	set := func(x, y int) {
		if x >= 0 && x < width && y >= 0 && y < height {
			out.Set(x, y, theme.Color(theme.ColorNamePrimary))
		}
	}

	// Draw chevron arrows centered on the divider
	const (
		arm       = 12
		thickness = 4
		offset    = 28
	)
	cy := height / 2
	for i := -arm; i <= arm; i++ {
		depth := i
		if depth < 0 {
			depth = -depth
		}
		for t := 0; t < thickness; t++ {
			set(splitX-offset+depth+t, cy+i) // Left chevron
			set(splitX+offset-depth-t, cy+i) // Right chevron
		}
	}

	return out
}
