// UI + tilemap smoke window for Gomag.
//
//	cd engine && CGO_ENABLED=0 mise exec -- go run ./examples/hello_ui
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gogpu/gg"
	"github.com/just-Bri/gomakeagame/engine/debug"
	"github.com/just-Bri/gomakeagame/engine/font"
	"github.com/just-Bri/gomakeagame/engine/gfx"
	"github.com/just-Bri/gomakeagame/engine/grid"
	"github.com/just-Bri/gomakeagame/engine/tilemap"
	"github.com/just-Bri/gomakeagame/engine/ui"
)

func main() {
	g, err := grid.New(16, 10)
	if err != nil {
		log.Fatal(err)
	}
	for x := 4; x < 8; x++ {
		for y := 3; y < 7; y++ {
			g.SetWalkable(x, y, false)
		}
	}
	layer := tilemap.FromWalkable(g,
		tilemap.RGB{0.22, 0.34, 0.22},
		tilemap.RGB{0.35, 0.28, 0.24},
	)

	app := gfx.New(gfx.Config{
		Title:      "Gomag — hello_ui",
		Width:      960,
		Height:     540,
		Continuous: true,
	})

	var fps debug.FPS
	var last time.Time
	cash := 100
	font.MustDefault()

	app.OnDraw(func(dc *gg.Context, w, h int) {
		now := time.Now()
		dt := 1.0 / 60.0
		if !last.IsZero() {
			dt = now.Sub(last).Seconds()
		}
		last = now
		fps.Tick(dt)

		lb := gfx.Fit(float64(w), float64(h), 960, 540)
		lb.Apply(dc)
		defer dc.Pop()

		dc.ClearWithColor(gg.RGB(0.09, 0.10, 0.12))

		tilemap.DrawOrtho(dc, layer, 36, 40, 80)
		tilemap.DrawOrthoGridLines(dc, layer.Width, layer.Height, 36, 40, 80, tilemap.RGB{0, 0, 0}, 0.25)

		ui.Panel(dc, 40, 20, 280, 48, ui.DefaultStyle())
		ui.Label(dc, "Maze Preview", 56, 50, 18, ui.RGB{0.95, 0.88, 0.55})

		ui.Panel(dc, 700, 20, 220, 90, ui.DefaultStyle())
		ui.Label(dc, fmt.Sprintf("Cash: %d", cash), 720, 50, 16, ui.RGB{0.9, 0.9, 0.8})
		ui.Bar(dc, 720, 70, 180, 14, 0.65, ui.RGB{0.3, 0.7, 0.35})
		ui.Label(dc, "Wave readiness", 720, 100, 12, ui.RGB{0.7, 0.7, 0.65})

		ov := debug.Overlay{}
		ov.Add("%s", fps.String())
		ov.Add("tiles %dx%d", layer.Width, layer.Height)
		ov.X, ov.Y = 40, 460
		ov.Draw(dc)
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
