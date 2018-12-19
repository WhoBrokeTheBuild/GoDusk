package main

import (
	"github.com/WhoBrokeTheBuild/GoDusk/dusk"
	_ "github.com/WhoBrokeTheBuild/GoDusk/dusk/obj"
	"github.com/WhoBrokeTheBuild/GoDusk/m32"
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type lightingShader struct {
	dusk.DefaultShader
	PointLight *lightEntity
}

func newLightingShader(pointLight *lightEntity) *lightingShader {
	s := &lightingShader{
		PointLight: pointLight,
	}
	s.InitFromFiles(
		"data/shaders/main.vs.glsl",
		"data/shaders/main.fs.glsl",
	)
	return s
}

func (s *lightingShader) Bind(ctx *dusk.RenderContext, data interface{}) {
	s.Shader.Bind(ctx, data)
	model := mgl32.Mat4{}
	if data != nil {
		model = data.(mgl32.Mat4)
	}

	mvp := ctx.Projection.
		Mul4(ctx.Camera.View).
		Mul4(model)

	gl.UniformMatrix4fv(s.UniformLocation("uProjection"), 1, false, &ctx.Projection[0])
	gl.UniformMatrix4fv(s.UniformLocation("uView"), 1, false, &ctx.Camera.View[0])
	gl.UniformMatrix4fv(s.UniformLocation("uModel"), 1, false, &model[0])
	gl.UniformMatrix4fv(s.UniformLocation("uMVP"), 1, false, &mvp[0])

	gl.Uniform3fv(s.UniformLocation("uPointLightPos"), 1, &s.PointLight.Transform().Position[0])
}

type flatShader struct {
	dusk.DefaultShader
	Color mgl32.Vec4
}

func newFlatShader(color mgl32.Vec4) *flatShader {
	s := &flatShader{
		Color: color,
	}
	s.InitFromFiles(
		"data/shaders/flat.vs.glsl",
		"data/shaders/flat.fs.glsl",
	)
	return s
}

func (s *flatShader) Bind(ctx *dusk.RenderContext, data interface{}) {
	s.Shader.Bind(ctx, data)
	model := mgl32.Mat4{}
	if data != nil {
		model = data.(mgl32.Mat4)
	}

	mvp := ctx.Projection.
		Mul4(ctx.Camera.View).
		Mul4(model)

	gl.UniformMatrix4fv(s.UniformLocation("uMVP"), 1, false, &mvp[0])

	gl.Uniform4fv(s.UniformLocation("uColor"), 1, &s.Color[0])
}

type lightEntity struct {
	dusk.Entity
	basePosition mgl32.Vec3
}

func newLightEntity(basePos mgl32.Vec3, layer dusk.ILayer) *lightEntity {
	e := &lightEntity{
		basePosition: basePos,
	}
	e.Init(layer)
	return e
}

func (e *lightEntity) Update(ctx *dusk.UpdateContext) {
	e.Entity.Update(ctx)

	e.Transform().Position[1] = e.basePosition[1] + m32.Sin(float32(ctx.ElapsedTime)*1.2)*3
}

type rotatingEntity struct {
	dusk.Entity
	rotInd int
}

func newRotatingEntity(rotInd int, layer dusk.ILayer) *rotatingEntity {
	e := &rotatingEntity{
		rotInd: rotInd,
	}
	e.Init(layer)
	return e
}

func (e *rotatingEntity) Update(ctx *dusk.UpdateContext) {
	e.Entity.Update(ctx)

	e.Transform().Rotation[e.rotInd] += ctx.DeltaTime * 0.01
	e.Transform().Rotation[e.rotInd] = float32(m32.Mod(e.Transform().Rotation[e.rotInd], m32.Pi*2.0))
}

func main() {
	opts := dusk.DefaultAppOptions()
	opts.Window.Title = "GoDusk - Mesh Lighting Demo"
	app, err := dusk.NewApp(opts)
	if err != nil {
		panic(err)
	}

	c := app.GetRenderContext().Camera
	c.SetPosition(mgl32.Vec3{8, 8, 8})
	c.SetLookAt(mgl32.Vec3{0, 1, 0})

	// Main Layer
	layer := dusk.NewLayer()
	app.AddLayer(layer)
	defer layer.Delete()

	// UI Layer
	ui, err := dusk.NewUILayer(app)
	app.AddLayer(ui)
	defer ui.Delete()

	fs := newFlatShader(mgl32.Vec4{1, 1, 0.9, 1})
	defer fs.Delete()

	l1 := newLightEntity(mgl32.Vec3{0, 3, 0}, layer)
	l1.Transform().Scale = mgl32.Vec3{0.2, 0.2, 0.2}
	defer l1.Delete()
	layer.AddEntity(l1)

	l1m, err := dusk.NewModelFromFile(l1, "data/models/uvsphere.obj")
	if err != nil {
		panic(err)
	}
	defer l1m.Delete()
	l1m.Shader = fs
	l1.AddComponent(l1m)

	s := newLightingShader(l1)
	defer s.Delete()

	r1 := newRotatingEntity(0, layer)
	r1.Transform().Position = mgl32.Vec3{-3, 1, -1}
	defer r1.Delete()
	layer.AddEntity(r1)

	r1m, err := dusk.NewModelFromFile(r1, "data/models/torus.obj")
	if err != nil {
		panic(err)
	}
	defer r1m.Delete()
	r1m.Shader = s
	r1.AddComponent(r1m)

	r2 := newRotatingEntity(2, layer)
	r2.Transform().Position = mgl32.Vec3{-5, 2, -2}
	defer r2.Delete()
	layer.AddEntity(r2)

	r2m, err := dusk.NewModelFromFile(r2, "data/models/torus.obj")
	if err != nil {
		panic(err)
	}
	defer r2m.Delete()
	r2m.Shader = s
	r2.AddComponent(r2m)

	r3 := newRotatingEntity(0, layer)
	r3.Transform().Position = mgl32.Vec3{-7, 3, -3}
	defer r3.Delete()
	layer.AddEntity(r3)

	r3m, err := dusk.NewModelFromFile(r3, "data/models/torus.obj")
	if err != nil {
		panic(err)
	}
	defer r3m.Delete()
	r3m.Shader = s
	r3.AddComponent(r3m)

	m1 := newRotatingEntity(1, layer)
	m1.Transform().Position = mgl32.Vec3{2, 0, -2}
	defer m1.Delete()
	layer.AddEntity(m1)

	m1m, err := dusk.NewModelFromFile(m1, "data/models/monkey.obj")
	if err != nil {
		panic(err)
	}
	defer m1m.Delete()
	m1m.Shader = s
	m1.AddComponent(m1m)

	m2 := dusk.NewEntity(layer)
	m2.Transform().Position = mgl32.Vec3{0, -2, 0}
	m2.Transform().Scale = mgl32.Vec3{10, 0.2, 10}
	defer m2.Delete()
	layer.AddEntity(m2)

	m2m, err := dusk.NewModelFromFile(m2, "data/models/cube.obj")
	if err != nil {
		panic(err)
	}
	defer m2m.Delete()
	m2m.Shader = s
	m2.AddComponent(m2m)

	m3 := newRotatingEntity(1, layer)
	m3.Transform().Position = mgl32.Vec3{0, 2, 5}
	m3.Transform().Rotation = mgl32.Vec3{1, 2, 3}
	defer m3.Delete()
	layer.AddEntity(m3)

	m3m, err := dusk.NewModelFromFile(m3, "data/models/cube.obj")
	if err != nil {
		panic(err)
	}
	defer m3m.Delete()
	m3m.Shader = s
	m3.AddComponent(m3m)

	app.Run()
}
