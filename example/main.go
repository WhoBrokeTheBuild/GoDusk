package main

//go:generate go-bindata -tags !release -debug -pkg $GOPACKAGE -o data.gen.go data/...
//go:generate go-bindata -tags release -pkg $GOPACKAGE -o data-release.gen.go data/...

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

	menu := dusk.NewUIText(fmt.Sprintf("GoDusk Example v%s %s", dusk.Version, Build), "data/fonts/default.ttf", 18.0, color.White)
	menu.SetPosition(mgl32.Vec2{10, 5})
	app.UI.AddElement(menu)

	fps := dusk.NewUIText("FPS 00", "data/fonts/default.ttf", 18.0, color.White)
	fps.SetPosition(mgl32.Vec2{float32(app.Window.Width) - 60, 5})
	app.UI.AddElement(fps)

	test, err := dusk.NewActor()
	if err != nil {
		panic(err)
	}

	mesh, err := dusk.NewMeshFromFile("data/models/teapot.fbx")
	if err != nil {
		panic(err)
	}
	test.AddMesh(mesh)

	cam := app.GetRenderContext().Camera
	horizontalAngle := float32(math.Pi * 1.25)
	verticalAngle := float32(math.Pi * -0.2)

	mPos := mgl32.Vec2{}
	mouseDown := false

	app.Window.RegisterMouseFunc(func(button dusk.MouseButton, action dusk.InputAction) {
		if button == dusk.MouseButtonLeft {
			if action == dusk.Press {
				mPos = app.Window.GetMousePos()
				mouseDown = true
			} else if action == dusk.Release {
				mouseDown = false
			}
		}
	})

	app.Window.RegisterMouseMoveFunc(func(pos mgl32.Vec2) {
		if mouseDown {
			delta := pos.Sub(mPos).Mul(0.01)
			horizontalAngle -= delta[0]
			verticalAngle -= delta[1]
			mPos = pos
		}
	})

	app.Window.RegisterKeyFunc(func(key dusk.Key, action dusk.InputAction) {
		// if action == dusk.Press {
		// 	switch key {
		// 	case dusk.KeyLeft:
		// 		turnSpeed = 0.1
		// 	case dusk.KeyRight:
		// 		turnSpeed = -0.1
		// 	case dusk.KeyUp:
		// 		moveSpeed = -1.0
		// 	case dusk.KeyDown:
		// 		moveSpeed = 1.0
		// 	}
		// } else if action == dusk.Release {
		// 	switch key {
		// 	case dusk.KeyLeft:
		// 		if turnSpeed > 0 {
		// 			turnSpeed = 0
		// 		}
		// 	case dusk.KeyRight:
		// 		if turnSpeed < 0 {
		// 			turnSpeed = 0
		// 		}
		// 	case dusk.KeyUp:
		// 		if moveSpeed < 0 {
		// 			moveSpeed = 0
		// 		}
		// 	case dusk.KeyDown:
		// 		if moveSpeed > 0 {
		// 			moveSpeed = 0
		// 		}
		// 	}
		// }
	})

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

		//cam.SetPosition()

		camDir := mgl32.Vec3{
			float32(math.Cos(float64(verticalAngle)) * math.Sin(float64(horizontalAngle))),
			float32(math.Sin(float64(verticalAngle))),
			float32(math.Cos(float64(verticalAngle)) * math.Cos(float64(horizontalAngle))),
		}
		cam.SetLookAt(cam.Position.Add(camDir))

		//test.Transform.Rotation[1] += ctx.DeltaTime * 0.01
		//test.Transform.Rotation[1] = float32(math.Mod(float64(test.Transform.Rotation[1]), math.Pi*2.0))
	})

	app.RegisterRenderFunc(func(ctx *dusk.RenderContext) {
		test.Render(ctx)
	})

	app.Run()
}
