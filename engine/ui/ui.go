// Package ui provides immediate-mode 2D HUD primitives for gg.
//
// Layout and game chrome stay in the game; this package draws reusable
// panels, bars, labels, and hit tests in virtual (letterboxed) coordinates.
package ui

import (
	"github.com/gogpu/gg"
	"github.com/just-Bri/gomakeagame/engine/font"
)

// RGB is a linear-ish 0..1 color (matches gg.SetRGB).
type RGB struct{ R, G, B float64 }

// Style controls panel chrome. Zero values get Classic-ish defaults.
type Style struct {
	Fill       RGB
	FillAlpha  float64
	Border     RGB
	BorderWidth float64
	Radius     float64
}

// DefaultStyle is a dark panel with warm gold trim (works for RPG HUD and TD overlays).
func DefaultStyle() Style {
	return Style{
		Fill:        RGB{0.08, 0.07, 0.06},
		FillAlpha:   0.82,
		Border:      RGB{0.72, 0.58, 0.22},
		BorderWidth: 1.5,
		Radius:      4,
	}
}

func (s Style) withDefaults() Style {
	d := DefaultStyle()
	if s.FillAlpha == 0 && s.Fill == (RGB{}) && s.Border == (RGB{}) {
		return d
	}
	if s.FillAlpha == 0 {
		s.FillAlpha = d.FillAlpha
	}
	if s.BorderWidth == 0 {
		s.BorderWidth = d.BorderWidth
	}
	if s.Radius == 0 {
		s.Radius = d.Radius
	}
	if s.Border == (RGB{}) {
		s.Border = d.Border
	}
	if s.Fill == (RGB{}) {
		s.Fill = d.Fill
	}
	return s
}

// Panel draws a rounded framed rectangle.
func Panel(dc *gg.Context, x, y, w, h float64, style Style) {
	if dc == nil {
		return
	}
	s := style.withDefaults()
	dc.SetRGBA(s.Fill.R, s.Fill.G, s.Fill.B, s.FillAlpha)
	dc.DrawRoundedRectangle(x, y, w, h, s.Radius)
	dc.Fill()
	dc.SetRGB(s.Border.R, s.Border.G, s.Border.B)
	dc.SetLineWidth(s.BorderWidth)
	dc.DrawRoundedRectangle(x, y, w, h, s.Radius)
	dc.Stroke()
}

// Bar draws a horizontal progress bar. frac is clamped to [0,1].
func Bar(dc *gg.Context, x, y, w, h, frac float64, fill RGB) {
	if dc == nil {
		return
	}
	if frac < 0 {
		frac = 0
	}
	if frac > 1 {
		frac = 1
	}
	dc.SetRGB(0.12, 0.12, 0.12)
	dc.DrawRectangle(x, y, w, h)
	dc.Fill()
	dc.SetRGB(fill.R, fill.G, fill.B)
	dc.DrawRectangle(x, y, w*frac, h)
	dc.Fill()
	dc.SetRGBA(1, 1, 1, 0.15)
	dc.DrawRectangle(x, y, w, h*0.35)
	dc.Fill()
}

// Label draws text with the default engine font at size.
func Label(dc *gg.Context, s string, x, y, size float64, color RGB) {
	if dc == nil || s == "" {
		return
	}
	font.Apply(dc, size)
	dc.SetRGB(color.R, color.G, color.B)
	dc.DrawString(s, x, y)
}

// LabelCentered draws text centered on (x,y).
func LabelCentered(dc *gg.Context, s string, x, y, size float64, color RGB) {
	if dc == nil || s == "" {
		return
	}
	font.Apply(dc, size)
	dc.SetRGB(color.R, color.G, color.B)
	dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
}

// Hit reports whether (px,py) lies inside the axis-aligned box.
func Hit(x, y, w, h, px, py float64) bool {
	return px >= x && py >= y && px < x+w && py < y+h
}

// Button draws a panel and returns whether it was pressed this frame.
// pressed should be true on the pointer-down edge while the cursor is over the button.
func Button(dc *gg.Context, x, y, w, h float64, label string, style Style, hovered, pressed bool) bool {
	s := style.withDefaults()
	if hovered {
		s.Fill = RGB{s.Fill.R + 0.06, s.Fill.G + 0.06, s.Fill.B + 0.06}
	}
	Panel(dc, x, y, w, h, s)
	LabelCentered(dc, label, x+w/2, y+h/2, 14, RGB{0.95, 0.9, 0.7})
	return pressed && hovered
}
