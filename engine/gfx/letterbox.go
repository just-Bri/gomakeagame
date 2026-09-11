package gfx

import (
	"math"

	"github.com/gogpu/gg"
)

// Letterbox maps a fixed virtual (playfield) size into a window while
// preserving aspect ratio, with black bars on the unused sides.
type Letterbox struct {
	Scale            float64
	OffsetX, OffsetY float64
	VirtW, VirtH     float64
}

// Fit computes a letterbox for virtual size (virtW×virtH) inside winW×winH.
// Degenerate sizes yield Scale=0 and zero offsets.
func Fit(winW, winH, virtW, virtH float64) Letterbox {
	if winW <= 0 || winH <= 0 || virtW <= 0 || virtH <= 0 {
		return Letterbox{VirtW: virtW, VirtH: virtH}
	}
	scale := math.Min(winW/virtW, winH/virtH)
	return Letterbox{
		Scale:   scale,
		OffsetX: (winW - virtW*scale) / 2,
		OffsetY: (winH - virtH*scale) / 2,
		VirtW:   virtW,
		VirtH:   virtH,
	}
}

// Apply pushes a transform so subsequent drawing uses virtual coordinates.
// Caller must dc.Pop() when finished (typically via defer).
func (l Letterbox) Apply(dc *gg.Context) {
	if dc == nil {
		return
	}
	dc.Push()
	dc.Translate(l.OffsetX, l.OffsetY)
	dc.Scale(l.Scale, l.Scale)
}

// ToVirtual converts window/screen pixels into virtual playfield coordinates.
func (l Letterbox) ToVirtual(sx, sy float64) (vx, vy float64) {
	if l.Scale == 0 {
		return 0, 0
	}
	return (sx - l.OffsetX) / l.Scale, (sy - l.OffsetY) / l.Scale
}

// ToScreen converts virtual playfield coordinates into window/screen pixels.
func (l Letterbox) ToScreen(vx, vy float64) (sx, sy float64) {
	return vx*l.Scale + l.OffsetX, vy*l.Scale + l.OffsetY
}
