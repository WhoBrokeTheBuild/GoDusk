package dusk

import "github.com/go-gl/mathgl/mgl32"

// Transform represents a position, rotation, and scale
type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}

// NewTransform creates a default, identity transformation
func NewTransform() *Transform {
	return &Transform{
		Position: mgl32.Vec3{},
		Rotation: mgl32.Vec3{},
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}

// GetMatrix returns the calculated 4x4 Matrix
func (t *Transform) GetMatrix() mgl32.Mat4 {
	return mgl32.Ident4().
		Mul4(mgl32.Translate3D(t.Position[0], t.Position[1], t.Position[2])).
		Mul4(mgl32.HomogRotate3DX(t.Rotation[0])).
		Mul4(mgl32.HomogRotate3DY(t.Rotation[1])).
		Mul4(mgl32.HomogRotate3DZ(t.Rotation[2])).
		Mul4(mgl32.Scale3D(t.Scale[0], t.Scale[1], t.Scale[2]))
}
