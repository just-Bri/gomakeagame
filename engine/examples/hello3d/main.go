// Minimal gogpu + g3d smoke window for Gomag.
//
//	cd engine && CGO_ENABLED=0 mise exec -- go run ./examples/hello3d
package main

import (
	"log"

	"github.com/gogpu/g3d"
	"github.com/just-Bri/gomakeagame/engine/gfx3d"
)

func main() {
	scene := g3d.NewScene()
	scene.SetBackground(g3d.RGB(0.10, 0.11, 0.14))

	ambient := g3d.NewAmbientLight(
		g3d.WithLightColor(g3d.White),
		g3d.WithLightIntensity(0.35),
	)
	ambientNode := g3d.NewNode()
	ambientNode.SetName("AmbientLight")
	ambientNode.SetUserData(ambient)
	scene.Add(ambientNode)

	sun := g3d.NewDirectionalLight(
		g3d.WithLightColor(g3d.White),
		g3d.WithLightIntensity(1.0),
	)
	sun.LightNode().SetRotation(g3d.Euler{
		X: g3d.Radians(-50),
		Y: g3d.Radians(35),
	})
	scene.Add(sun.LightNode())

	ground := g3d.NewMesh(
		g3d.NewPlaneGeometry(8, 8),
		g3d.NewStandardMaterial(
			g3d.WithColor(g3d.RGB(0.22, 0.28, 0.20)),
			g3d.WithRoughness(0.95),
		),
	)
	ground.MeshNode().SetRotation(g3d.Euler{X: g3d.Radians(-90)})
	scene.Add(ground.MeshNode())

	cube := g3d.NewMesh(
		g3d.NewBoxGeometry(1, 1, 1),
		g3d.NewStandardMaterial(
			g3d.WithColor(g3d.RGB(0.45, 0.70, 0.95)),
			g3d.WithMetallic(0.25),
			g3d.WithRoughness(0.55),
		),
	)
	cube.MeshNode().SetPosition(g3d.Vec3{Y: 0.5})
	scene.Add(cube.MeshNode())

	camera := g3d.NewPerspectiveCamera(60, 1280.0/720.0, 0.1, 200)
	camera.CameraNode().SetPosition(g3d.Vec3{X: 3.2, Y: 2.4, Z: 3.2})
	camera.CameraNode().LookAt(g3d.Vec3{Y: 0.4})

	app := gfx3d.New(gfx3d.Config{
		Title:      "Gomag — hello3d",
		Width:      960,
		Height:     540,
		Continuous: true,
	})
	app.SetScene(scene, camera)
	app.OnUpdate(func(dt float64) {
		r := cube.MeshNode().Rotation
		r.Y += float32(dt)
		cube.MeshNode().SetRotation(r)
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
