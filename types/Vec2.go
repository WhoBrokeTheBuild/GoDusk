package types

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
)

type Vec2 mgl32.Vec2
type Vec3 mgl32.Vec3
type Vec4 mgl32.Vec4

type Vec2d mgl64.Vec2
type Vec3d mgl64.Vec3
type Vec4d mgl64.Vec4

type Vec2i [2]int

func (v Vec2i) Elem() (int, int) {
	return v[0], v[1]
}

func (v Vec2i) X() int {
	return v[0]
}

func (v Vec2i) Y() int {
	return v[1]
}

type Vec3i [3]int

func (v Vec3i) Elem() (int, int, int) {
	return v[0], v[1], v[2]
}

func (v Vec3i) X() int {
	return v[0]
}

func (v Vec3i) Y() int {
	return v[1]
}

func (v Vec3i) Z() int {
	return v[2]
}

type Vec4i [4]int

func (v Vec4i) Elem() (int, int, int, int) {
	return v[0], v[1], v[2], v[3]
}

func (v Vec4i) X() int {
	return v[0]
}

func (v Vec4i) Y() int {
	return v[1]
}

func (v Vec4i) Z() int {
	return v[2]
}

func (v Vec4i) W() int {
	return v[3]
}
