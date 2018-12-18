package main

//go:generate go-bindata -tags !release -debug -pkg $GOPACKAGE -o data.gen.go data/...
//go:generate go-bindata -tags release -pkg $GOPACKAGE -o data-release.gen.go data/...

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/fbx"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/obj"
	"github.com/WhoBrokeTheBuild/GoDusk/m32"
)

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
	opts.Window.Title = "GoDusk - Mesh Loading Demo"
	app, err := dusk.NewApp(opts)
	if err != nil {
		panic(err)
	}

	fbxActor := newDemoActor()
	defer fbxActor.Delete()
	fbxActor.Transform().Position = mgl32.Vec3{2, 0, 0}
	app.Scene.AddActor(fbxActor)

	fbxMesh, err := dusk.NewMeshFromFile("data/models/teapot.fbx")
	if err != nil {
		panic(err)
	}
	defer fbxMesh.Delete()
	fbxActor.AddMesh(fbxMesh)

	fbxLabel := dusk.NewUIText("teapot.fbx", "data/fonts/default.ttf", 26.0, color.Black)
	fbxLabel.SetPosition(mgl32.Vec2{700, 500})
	app.UI.AddElement(fbxLabel)

	objActor := newDemoActor()
	defer objActor.Delete()
	objActor.Transform().Position = mgl32.Vec3{0, 0, 2}
	app.Scene.AddActor(objActor)

	objMesh, err := dusk.NewMeshFromFile("data/models/teapot.obj")
	if err != nil {
		panic(err)
	}
	defer objMesh.Delete()
	objActor.AddMesh(objMesh)

	objLabel := dusk.NewUIText("teapot.obj", "data/fonts/default.ttf", 26.0, color.Black)
	objLabel.SetPosition(mgl32.Vec2{200, 500})
	app.UI.AddElement(objLabel)

	cam := app.GetRenderContext().Camera
	cam.SetPosition(mgl32.Vec3{4, 4, 4})
	cam.SetLookAt(mgl32.Vec3{1, 0, 1})

	app.Run()
}
