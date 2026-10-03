package project

import (
	"sort"

	"github.com/gogpu/gg"
)

// DrawItem is one painter's-algorithm entry for soft-3D scenes.
type DrawItem struct {
	Depth float64
	Draw  func(dc *gg.Context)
}

// DrawSorted paints items back-to-front by Depth (smaller depth first).
func DrawSorted(dc *gg.Context, items []DrawItem) {
	if dc == nil || len(items) == 0 {
		return
	}
	sort.SliceStable(items, func(i, j int) bool { return items[i].Depth < items[j].Depth })
	for _, it := range items {
		if it.Draw != nil {
			it.Draw(dc)
		}
	}
}

// DrawBox renders an axis-aligned placeholder box through an isometric projector.
func DrawBox(dc *gg.Context, iso Iso, x, y, z, w, h, d, r, g, b float64) {
	if dc == nil {
		return
	}
	corners := BoxCorners(x, y, z, w, h, d)
	screen := [8][2]float64{}
	for i, c := range corners {
		sx, sy := iso.WorldToScreen(c[0], c[1], c[2])
		screen[i] = [2]float64{sx, sy}
	}
	type face struct {
		idx   int
		depth float64
	}
	faces := make([]face, 0, 6)
	for fi, idx := range FaceIndices {
		avg := 0.0
		for _, ci := range idx {
			avg += iso.DepthKey(corners[ci][0], corners[ci][1], corners[ci][2])
		}
		faces = append(faces, face{idx: fi, depth: avg / 4})
	}
	sort.SliceStable(faces, func(i, j int) bool { return faces[i].depth < faces[j].depth })
	for _, f := range faces {
		idx := FaceIndices[f.idx]
		shade := FaceShade(f.idx)
		dc.SetRGB(Clamp01(r*shade), Clamp01(g*shade), Clamp01(b*shade))
		dc.MoveTo(screen[idx[0]][0], screen[idx[0]][1])
		dc.LineTo(screen[idx[1]][0], screen[idx[1]][1])
		dc.LineTo(screen[idx[2]][0], screen[idx[2]][1])
		dc.LineTo(screen[idx[3]][0], screen[idx[3]][1])
		dc.ClosePath()
		dc.Fill()
		dc.SetRGBA(0, 0, 0, 0.25)
		dc.SetLineWidth(1)
		dc.MoveTo(screen[idx[0]][0], screen[idx[0]][1])
		dc.LineTo(screen[idx[1]][0], screen[idx[1]][1])
		dc.LineTo(screen[idx[2]][0], screen[idx[2]][1])
		dc.LineTo(screen[idx[3]][0], screen[idx[3]][1])
		dc.ClosePath()
		dc.Stroke()
	}
}
