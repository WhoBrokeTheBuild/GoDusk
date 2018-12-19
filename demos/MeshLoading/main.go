package main

import (
	"image/color"

	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/fbx"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/obj"
	"github.com/WhoBrokeTheBuild/GoDusk/m32"
)

type demoEntity struct {
	dusk.Entity
}

func newDemoEntity(layer dusk.ILayer) *demoEntity {
	e := &demoEntity{}
	e.Init(layer)
	return e
}

func (e *demoEntity) Update(ctx *dusk.UpdateContext) {
	e.Entity.Update(ctx)

	e.Transform().Rotation[1] += ctx.DeltaTime * 0.01
	e.Transform().Rotation[1] = float32(m32.Mod(e.Transform().Rotation[1], m32.Pi*2.0))
}

func main() {
	opts := dusk.DefaultAppOptions()
	opts.Window.Title = "GoDusk - Mesh Loading Demo"
	app, err := dusk.NewApp(opts)
	if err != nil {
		panic(err)
	}

	layer := dusk.NewLayer()
	app.AddLayer(layer)
	defer layer.Delete()

	ui, err := dusk.NewUILayer(app)
	app.AddLayer(ui)
	defer ui.Delete()

	fbxEntity := newDemoEntity(layer)
	fbxEntity.Transform().Position = mgl32.Vec3{2, 0, 0}
	layer.AddEntity(fbxEntity)
	defer fbxEntity.Delete()

	fbxModel, err := dusk.NewModelFromFile(fbxEntity, "data/models/teapot.fbx")
	if err != nil {
		panic(err)
	}
	defer fbxModel.Delete()
	fbxEntity.AddComponent(fbxModel)

	fbxLabel := dusk.NewUIText(ui, "teapot.fbx", "data/fonts/default.ttf", 26.0, color.Black)
	fbxLabel.SetPosition(mgl32.Vec2{700, 500})
	ui.AddEntity(fbxLabel)

	objEntity := newDemoEntity(layer)
	objEntity.Transform().Position = mgl32.Vec3{0, 0, 2}
	layer.AddEntity(objEntity)
	defer objEntity.Delete()

	objModel, err := dusk.NewModelFromFile(objEntity, "data/models/teapot.obj")
	if err != nil {
		panic(err)
	}
	defer objModel.Delete()
	objEntity.AddComponent(objModel)

	objLabel := dusk.NewUIText(ui, "teapot.obj", "data/fonts/default.ttf", 26.0, color.Black)
	objLabel.SetPosition(mgl32.Vec2{200, 500})
	ui.AddEntity(objLabel)

	cam := app.GetRenderContext().Camera
	cam.SetPosition(mgl32.Vec3{4, 4, 4})
	cam.SetLookAt(mgl32.Vec3{1, 0, 1})

	app.Run()
}
