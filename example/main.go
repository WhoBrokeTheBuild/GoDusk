package main

//go:generate go-bindata -tags !release -debug -pkg $GOPACKAGE -o data.gen.go data/...
//go:generate go-bindata -tags release -pkg $GOPACKAGE -o data-release.gen.go data/...

import (
	"fmt"
	"image/color"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/fbx"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/obj"
	"github.com/WhoBrokeTheBuild/GoDusk/m32"
)

// Build is a build identifier, generally the Git Short
var Build = ""

type demoActor struct {
	dusk.Actor
}

func newDemoActor() *demoActor {
	a := &demoActor{}
	a.Init()
	return a
}

func (a *demoActor) Update(ctx *dusk.UpdateContext) {
	a.Actor.Update(ctx)

	a.Transform().Rotation[1] += ctx.DeltaTime * 0.01
	a.Transform().Rotation[1] = float32(m32.Mod(a.Transform().Rotation[1], m32.Pi*2.0))
}

func main() {
	dusk.RegisterFunc(Asset)

	opts := dusk.DefaultAppOptions()
	opts.Window.Icons = []string{"data/icons/icon_64.png", "data/icons/icon_32.png"}
	app, err := dusk.NewApp(opts)
	if err != nil {
		panic(err)
	}

	app.UI.AddElement(dusk.NewUIImageFromFile("data/ui/menubar.png"))

	menu := dusk.NewUIText(fmt.Sprintf("GoDusk Example v%s %s", dusk.Version, Build), "data/fonts/default.ttf", 18.0, color.White)
	menu.SetPosition(mgl32.Vec2{10, 5})
	app.UI.AddElement(menu)

	fps := dusk.NewUIText("FPS 00", "data/fonts/default.ttf", 18.0, color.White)
	fps.SetPosition(mgl32.Vec2{float32(app.Window.Width) - 60, 5})
	app.UI.AddElement(fps)

	actor := newDemoActor()
	defer actor.Delete()
	app.Scene.AddActor(actor)

	mesh, err := dusk.NewMeshFromFile("data/models/teapot.obj")
	if err != nil {
		panic(err)
	}
	defer mesh.Delete()
	actor.AddMesh(mesh)

	//mPos := mgl32.Vec2{}
	mouseDown := false

	app.Window.RegisterMouseFunc(func(button dusk.MouseButton, action dusk.InputAction) {
		if button == dusk.MouseButtonLeft {
			if action == dusk.Press {
				//mPos = app.Window.GetMousePos()
				mouseDown = true
			} else if action == dusk.Release {
				mouseDown = false
			}
		}
	})

	app.Window.RegisterMouseMoveFunc(func(pos mgl32.Vec2) {
		if mouseDown {
			//delta := pos.Sub(mPos).Mul(0.01)
			//horizontalAngle -= delta[0]
			//verticalAngle -= delta[1]
			//mPos = pos
		}
	})

	//app.Window.RegisterKeyFunc(func(key dusk.Key, action dusk.InputAction) { })

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

		//camDir := mgl32.Vec3{
		//	float32(m32.Cos(verticalAngle) * m32.Sin(horizontalAngle)),
		//	float32(m32.Sin(verticalAngle)),
		//	float32(m32.Cos(verticalAngle) * m32.Cos(horizontalAngle)),
		//}
		//cam.SetLookAt(cam.Position.Add(camDir))
	})

	app.Run()
}
