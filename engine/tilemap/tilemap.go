// Package tilemap draws integer tile layers for ortho and isometric views.
//
// Walkability stays in engine/grid; this package is presentation only.
package tilemap

import (
	"github.com/gogpu/gg"
	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/project"
)

// RGB is a 0..1 color.
type RGB struct{ R, G, B float64 }

// Layer is a dense W×H map of palette indices (row-major).
// Index 0 is often "empty" / skipped when SkipZero is set.
type Layer struct {
	Width, Height int
	Tiles         []int
	Palette       []RGB
	SkipZero      bool
}

// NewLayer allocates a zeroed layer.
func NewLayer(w, h int, palette []RGB) Layer {
	return Layer{
		Width:   w,
		Height:  h,
		Tiles:   make([]int, w*h),
		Palette: palette,
	}
}

// At returns the tile id at (x,y), or 0 if out of bounds.
func (l Layer) At(x, y int) int {
	if x < 0 || y < 0 || x >= l.Width || y >= l.Height {
		return 0
	}
	return l.Tiles[y*l.Width+x]
}

// Set writes a tile id. Out of bounds is a no-op.
func (l Layer) Set(x, y, id int) {
	if x < 0 || y < 0 || x >= l.Width || y >= l.Height {
		return
	}
	l.Tiles[y*l.Width+x] = id
}

// FromWalkable builds a 2-color layer from a grid (1=walkable, 2=blocked).
func FromWalkable(g *grid.Grid, walk, blocked RGB) Layer {
	if g == nil {
		return Layer{}
	}
	l := NewLayer(g.Width, g.Height, []RGB{{}, walk, blocked})
	for y := 0; y < g.Height; y++ {
		for x := 0; x < g.Width; x++ {
			if g.Walkable(x, y) {
				l.Set(x, y, 1)
			} else {
				l.Set(x, y, 2)
			}
		}
	}
	return l
}

// DrawOrtho paints axis-aligned tiles. origin is the top-left of cell (0,0).
func DrawOrtho(dc *gg.Context, layer Layer, tilePx, originX, originY float64) {
	if dc == nil || layer.Width <= 0 || layer.Height <= 0 || tilePx <= 0 {
		return
	}
	for y := 0; y < layer.Height; y++ {
		for x := 0; x < layer.Width; x++ {
			id := layer.At(x, y)
			if layer.SkipZero && id == 0 {
				continue
			}
			c := colorOf(layer, id)
			dc.SetRGB(c.R, c.G, c.B)
			dc.DrawRectangle(originX+float64(x)*tilePx, originY+float64(y)*tilePx, tilePx, tilePx)
			dc.Fill()
		}
	}
}

// DrawOrthoGridLines strokes cell borders (debug / editor).
func DrawOrthoGridLines(dc *gg.Context, w, h int, tilePx, originX, originY float64, line RGB, alpha float64) {
	if dc == nil || w <= 0 || h <= 0 {
		return
	}
	dc.SetRGBA(line.R, line.G, line.B, alpha)
	dc.SetLineWidth(1)
	for x := 0; x <= w; x++ {
		px := originX + float64(x)*tilePx
		dc.DrawLine(px, originY, px, originY+float64(h)*tilePx)
		dc.Stroke()
	}
	for y := 0; y <= h; y++ {
		py := originY + float64(y)*tilePx
		dc.DrawLine(originX, py, originX+float64(w)*tilePx, py)
		dc.Stroke()
	}
}

// DrawIsoDiamonds paints each tile as an isometric diamond on Y=0.
// cellSize is the world-unit size of one tile (often 1).
func DrawIsoDiamonds(dc *gg.Context, layer Layer, iso project.Iso, cellSize float64) {
	if dc == nil || layer.Width <= 0 || layer.Height <= 0 {
		return
	}
	if cellSize <= 0 {
		cellSize = 1
	}
	for z := 0; z < layer.Height; z++ {
		for x := 0; x < layer.Width; x++ {
			id := layer.At(x, z)
			if layer.SkipZero && id == 0 {
				continue
			}
			c := colorOf(layer, id)
			wx := (float64(x) + 0.5) * cellSize
			wz := (float64(z) + 0.5) * cellSize
			hs := cellSize * 0.5
			x0, y0 := iso.WorldToScreen(wx-hs, 0, wz)
			x1, y1 := iso.WorldToScreen(wx, 0, wz-hs)
			x2, y2 := iso.WorldToScreen(wx+hs, 0, wz)
			x3, y3 := iso.WorldToScreen(wx, 0, wz+hs)
			dc.SetRGB(c.R, c.G, c.B)
			dc.MoveTo(x0, y0)
			dc.LineTo(x1, y1)
			dc.LineTo(x2, y2)
			dc.LineTo(x3, y3)
			dc.ClosePath()
			dc.Fill()
		}
	}
}

func colorOf(layer Layer, id int) RGB {
	if id >= 0 && id < len(layer.Palette) {
		return layer.Palette[id]
	}
	return RGB{0.5, 0.5, 0.5}
}
