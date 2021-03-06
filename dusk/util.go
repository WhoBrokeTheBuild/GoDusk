package dusk

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/WhoBrokeTheBuild/GoDusk/m32"
	"github.com/go-gl/mathgl/mgl32"
)

func get2DMeshData(dst, src mgl32.Vec4) *MeshData {
	return &MeshData{
		Vertices: []mgl32.Vec3{
			mgl32.Vec3{dst[2], dst[1], 0},
			mgl32.Vec3{dst[2], dst[3], 0},
			mgl32.Vec3{dst[0], dst[3], 0},
			mgl32.Vec3{dst[2], dst[1], 0},
			mgl32.Vec3{dst[0], dst[3], 0},
			mgl32.Vec3{dst[0], dst[1], 0},
		},
		TexCoords: []mgl32.Vec2{
			mgl32.Vec2{src[2], src[3]},
			mgl32.Vec2{src[2], src[1]},
			mgl32.Vec2{src[0], src[1]},
			mgl32.Vec2{src[2], src[3]},
			mgl32.Vec2{src[0], src[1]},
			mgl32.Vec2{src[0], src[3]},
		},
	}
}

func new2DMesh(dst, src mgl32.Vec4) (*Mesh, error) {
	mesh, err := NewMeshFromData(get2DMeshData(dst, src))
	if err != nil {
		return nil, err
	}
	return mesh, err
}

func update2DMesh(mesh *Mesh, dst, src mgl32.Vec4) error {
	err := mesh.UpdateData(get2DMeshData(dst, src))
	if err != nil {
		return err
	}
	return err
}

func addLineNumbers(code string) string {
	lines := strings.Split(code, "\n")
	w := int(m32.Log10(float32(len(lines)))) + 1
	f := "%" + strconv.Itoa(w) + "v: %v"
	for i := range lines {
		lines[i] = fmt.Sprintf(f, i+1, lines[i])
	}
	return strings.Join(lines, "\n")
}

// DistanceSquared returns the squared distance between two points
func DistanceSquared(p1, p2 mgl32.Vec3) float32 {
	tmp := p2.Sub(p1)
	return tmp.Dot(tmp)
}

// Distance returns the distance between two points
func Distance(p1, p2 mgl32.Vec3) float32 {
	return m32.Sqrt(DistanceSquared(p1, p2))
}
