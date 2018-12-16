package dusk

import "github.com/go-gl/mathgl/mgl32"

// RenderContext is a context of view and shader data
type RenderContext struct {
	Projection mgl32.Mat4
	Camera     *Camera
}
