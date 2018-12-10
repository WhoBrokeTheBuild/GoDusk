package main

//go:generate go-bindata -tags debug -debug -pkg $GOPACKAGE -o data-debug.gen.go data/...
//go:generate go-bindata -tags !debug -pkg $GOPACKAGE -o data.gen.go data/...

import (
	"fmt"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/context"
	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	"github.com/WhoBrokeTheBuild/GoDusk/load"
	"github.com/WhoBrokeTheBuild/GoDusk/ui"
)

var GIT_SHORT = ""

func main() {
	load.RegisterFunc(Asset)

	opts := dusk.DefaultAppOptions()
	opts.Window.Icons = []string{"data/icons/icon_64.png", "data/icons/icon_32.png"}
	app, err := dusk.NewApp(opts)
	if err != nil {
		panic(err)
	}

	app.UI.AddComponent(ui.NewImageFromFile("data/ui/menubar.png"))

	menu := ui.NewText(fmt.Sprintf("GoDusk Example v%s %s", dusk.Version, GIT_SHORT), "data/ui/default.ttf", 18.0, color.White)
	menu.SetPosition(mgl32.Vec2{10, 5})
	app.UI.AddComponent(menu)

	fps := ui.NewText("FPS 00", "data/ui/default.ttf", 18.0, color.White)
	fps.SetPosition(mgl32.Vec2{float32(app.Window.Width) - 60, 5})
	app.UI.AddComponent(fps)

	box, err := dusk.NewModelFromFile("data/models/crate/crate.obj")
	if err != nil {
		panic(err)
	}
	defer box.Delete()

	lastFPS := 0
	app.RegisterUpdateFunc(func(ctx *context.Update) {
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
	})

	app.RegisterRenderFunc(func(ctx *context.Render) {
		box.Render(ctx)
	})

	app.Run()
}
