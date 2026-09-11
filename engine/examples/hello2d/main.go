// Minimal gogpu + gg smoke window for Gomag.
//
//	cd engine && CGO_ENABLED=0 mise exec -- go run ./examples/hello2d
package main

import (
	"log"

	"github.com/gogpu/gg"
	"github.com/just-Bri/gomakeagame/engine/gfx"
)

func main() {
	app := gfx.New(gfx.Config{
		Title:      "Gomag — hello2d",
		Width:      960,
		Height:     540,
		Continuous: true,
	})

	app.OnDraw(func(dc *gg.Context, w, h int) {
		dc.ClearWithColor(gg.RGB(0.08, 0.09, 0.12))

		// Soft playfield panel.
		dc.SetRGB(0.14, 0.16, 0.20)
		margin := 40.0
		dc.DrawRoundedRectangle(margin, margin, float64(w)-2*margin, float64(h)-2*margin, 12)
		dc.Fill()

		// Accent circle (gg SDF path when GPU accel is on).
		cx, cy := float64(w)/2, float64(h)/2
		dc.SetRGB(0.35, 0.72, 0.55)
		dc.DrawCircle(cx, cy, 64)
		dc.Fill()

		dc.SetRGB(0.85, 0.88, 0.92)
		dc.DrawCircle(cx, cy, 64)
		dc.SetLineWidth(3)
		dc.Stroke()
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
