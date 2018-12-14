package main

//go:generate go-bindata -tags debug -debug -pkg $GOPACKAGE -o data-debug.gen.go data/...
//go:generate go-bindata -tags !debug -pkg $GOPACKAGE -o data.gen.go data/...

import (
	"fmt"
	"image/color"
	"math"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/fbx"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/obj"
)

// Build is a build identifier, generally the Git Short
var Build = ""

func main() {
	dusk.RegisterFunc(Asset)

	opts := dusk.DefaultAppOptions()
	opts.Window.Icons = []string{"data/icons/icon_64.png", "data/icons/icon_32.png"}
	app, err := dusk.NewApp(opts)
	if err != nil {
		panic(err)
	}

	app.UI.AddElement(dusk.NewUIImageFromFile("data/ui/menubar.png"))

	menu := dusk.NewUIText(fmt.Sprintf("GoDusk Example v%s %s", dusk.Version, Build), "data/ui/default.ttf", 18.0, color.White)
	menu.SetPosition(mgl32.Vec2{10, 5})
	app.UI.AddElement(menu)

	fps := dusk.NewUIText("FPS 00", "data/ui/default.ttf", 18.0, color.White)
	fps.SetPosition(mgl32.Vec2{float32(app.Window.Width) - 60, 5})
	app.UI.AddElement(fps)

	test, err := dusk.NewActor()
	if err != nil {
		panic(err)
	}

	mesh, err := dusk.NewMeshFromFile("data/models/teapot.obj")
	if err != nil {
		panic(err)
	}
	test.AddMesh(mesh)

	lastFPS := 0
	app.RegisterUpdateFunc(func(ctx *dusk.UpdateContext) {
		if lastFPS != ctx.FPS {
			lastFPS = ctx.FPS

			if ctx.FPS < 30 {
				fps.Color = color.RGBA{255, 0, 0, 255}
			} else if ctx.FPS < 60 {
				fps.Color = color.RGBA{255, 255, 0, 255}
			} else {
				fps.Color = color.RGBA{0, 255, 0, 255}
			}
			fps.SetText(fmt.Sprintf("FPS %d", ctx.FPS))
		}

		test.Transform.Rotation[1] += ctx.DeltaTime * 0.01
		test.Transform.Rotation[1] = float32(math.Mod(float64(test.Transform.Rotation[1]), math.Pi*2.0))
	})

	app.RegisterRenderFunc(func(ctx *dusk.RenderContext) {
		test.Render(ctx)
	})

	app.Run()
}
