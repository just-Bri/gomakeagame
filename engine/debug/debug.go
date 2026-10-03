// Package debug draws lightweight on-screen diagnostics.
package debug

import (
	"fmt"
	"time"

	"github.com/gogpu/gg"
	"github.com/just-Bri/gomakeagame/engine/font"
	"github.com/just-Bri/gomakeagame/engine/ui"
)

// Overlay is a stack of diagnostic lines (FPS, phase, entity counts, …).
type Overlay struct {
	Lines []string
	X, Y  float64 // top-left; zero → 8,8
}

// Add appends a formatted line.
func (o *Overlay) Add(format string, args ...any) {
	o.Lines = append(o.Lines, fmt.Sprintf(format, args...))
}

// Clear removes all lines.
func (o *Overlay) Clear() { o.Lines = o.Lines[:0] }

// Draw renders the overlay in the current coordinate space.
func (o Overlay) Draw(dc *gg.Context) {
	if dc == nil || len(o.Lines) == 0 {
		return
	}
	x, y := o.X, o.Y
	if x == 0 && y == 0 {
		x, y = 8, 8
	}
	h := 12.0 + float64(len(o.Lines))*16
	w := 220.0
	for _, line := range o.Lines {
		if float64(len(line))*7+24 > w {
			w = float64(len(line))*7 + 24
		}
	}
	ui.Panel(dc, x, y, w, h, ui.Style{
		Fill: ui.RGB{0.02, 0.02, 0.03}, FillAlpha: 0.7,
		Border: ui.RGB{0.4, 0.4, 0.45}, BorderWidth: 1, Radius: 3,
	})
	font.Apply(dc, 12)
	ty := y + 16
	for _, line := range o.Lines {
		dc.SetRGB(0.75, 0.9, 0.75)
		dc.DrawString(line, x+8, ty)
		ty += 16
	}
}

// FPS tracks a smoothed frames-per-second estimate.
type FPS struct {
	Value     float64
	last      time.Time
	frameCount int
	accum     float64
}

// Tick updates the FPS estimate; call once per frame with dt seconds.
func (f *FPS) Tick(dt float64) {
	if f == nil {
		return
	}
	if dt <= 0 {
		return
	}
	f.frameCount++
	f.accum += dt
	if f.accum >= 0.5 {
		f.Value = float64(f.frameCount) / f.accum
		f.frameCount = 0
		f.accum = 0
	}
}

// String returns a short "FPS 60" label.
func (f FPS) String() string {
	return fmt.Sprintf("FPS %.0f", f.Value)
}
