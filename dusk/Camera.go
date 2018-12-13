package dusk

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	Position mgl32.Vec3
	LookAt   mgl32.Vec3
	View     mgl32.Mat4
}

func NewCamera(pos, lookAt mgl32.Vec3) *Camera {
	c := &Camera{
		Position: pos,
		LookAt:   lookAt,
	}
	c.calcView()

	return c
}

func (c *Camera) SetPosition(pos mgl32.Vec3) {
	c.Position = pos
	c.calcView()
}

func (c *Camera) SetLookAt(lookAt mgl32.Vec3) {
	c.LookAt = lookAt
	c.calcView()
}

func (c *Camera) calcView() {
	c.View = mgl32.LookAtV(c.Position, c.LookAt, mgl32.Vec3{0, 1, 0})
}
