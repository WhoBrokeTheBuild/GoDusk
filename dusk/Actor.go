package dusk

import (
	"github.com/go-gl/gl/v4.1-core/gl"

	"github.com/go-gl/mathgl/mgl32"
)

type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}

func NewTransform() Transform {
	return Transform{
		Position: mgl32.Vec3{},
		Rotation: mgl32.Vec3{},
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}

func (t Transform) GetMatrix() mgl32.Mat4 {
	return mgl32.Ident4().
		Mul4(mgl32.Translate3D(t.Position[0], t.Position[1], t.Position[2])).
		Mul4(mgl32.HomogRotate3DX(t.Rotation[0])).
		Mul4(mgl32.HomogRotate3DY(t.Rotation[1])).
		Mul4(mgl32.HomogRotate3DZ(t.Rotation[2])).
		Mul4(mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2]))
}

// Actor is an object with a Transform
type Actor struct {
	Transform Transform
	Meshes    []*Mesh
}

// NewActor returns a new Actor
func NewActor() (*Actor, error) {
	a := &Actor{
		Transform: NewTransform(),
	}

	return a, nil
}

// AddMesh adds one or more meshes to the Actor
func (a *Actor) AddMesh(mesh ...*Mesh) {
	a.Meshes = append(a.Meshes, mesh...)
}

// Delete frees all resources owned by the Actor
func (a *Actor) Delete() {
}

func (a *Actor) Update(ctx *UpdateContext) {

}

// Render renders a Actor to the screen
func (a *Actor) Render(ctx *RenderContext) {
	s := ctx.Shader
	s.Bind()

	model := a.Transform.GetMatrix()

	gl.UniformMatrix4fv(s.GetUniformLocation("uProjection"), 1, false, &ctx.Projection[0])
	gl.UniformMatrix4fv(s.GetUniformLocation("uView"), 1, false, &ctx.Camera.View[0])
	gl.UniformMatrix4fv(s.GetUniformLocation("uModel"), 1, false, &model[0])
	gl.Uniform4fv(s.GetUniformLocation("uCamera"), 1, &ctx.Camera.Position[0])

	for _, m := range a.Meshes {
		m.Render(s)
	}
}
