package dusk

import (
	"C"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)
import (
	"fmt"
	"path/filepath"

	"github.com/WhoBrokeTheBuild/GoDusk/dusk/obj"
)

// Mesh represents a set of OpenGL Vertex Array Objects and Material data
type Mesh struct {
	groups []*meshGroup
}

type meshGroup struct {
	material *Material
	vao      uint32
	vbo      uint32
	size     int
	count    int32
}

// MeshData is the intermediate data format for loading Meshes from Memory
type MeshData struct {
	Material *Material

	Vertices  []mgl32.Vec3
	Normals   []mgl32.Vec3
	TexCoords []mgl32.Vec2
}

// NewMeshFromFile returns a new Mesh from the given file
func NewMeshFromFile(filename string) (*Mesh, error) {
	m := &Mesh{
		groups: []*meshGroup{},
	}

	err := m.LoadFromFile(filename)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, err
}

// NewMeshFromData returns a new Mesh from the given MeshData
func NewMeshFromData(data ...*MeshData) (*Mesh, error) {
	m := &Mesh{
		groups: []*meshGroup{},
	}

	err := m.LoadFromData(data...)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, err
}

// Delete frees all resources owned by the Mesh
func (m *Mesh) Delete() {
	for _, g := range m.groups {
		if g.material != nil {
			g.material.Delete()
		}
		if g.vbo != InvalidID {
			gl.DeleteBuffers(1, &g.vbo)
			g.vbo = InvalidID
		}
		if g.vao != InvalidID {
			gl.DeleteVertexArrays(1, &g.vao)
			g.vao = InvalidID
		}
	}
	m.groups = nil
}

// LoadFromFile loads a mesh from a given file
func (m *Mesh) LoadFromFile(filename string) error {
	filename = filepath.Clean(filename)
	m.Delete()

	Loadf("asset.Mesh [%v]", filename)
	r := obj.NewReaderEx(filename, obj.LoadFunc(Load))
	objs, err := r.Read()
	if err != nil {
		return err
	}

	if len(objs) == 0 {
		return fmt.Errorf("No objects loaded from [%v]", filename)
	}

	data := []*MeshData{}
	for _, o := range objs {
		mat, err := NewMaterialFromData(&MaterialData{
			Ambient:     mgl32.Vec4{o.Material.Ambient[0], o.Material.Ambient[1], o.Material.Ambient[2], 1},
			Diffuse:     mgl32.Vec4{o.Material.Diffuse[0], o.Material.Diffuse[1], o.Material.Diffuse[2], 1},
			Specular:    mgl32.Vec4{o.Material.Specular[0], o.Material.Specular[1], o.Material.Specular[2], 1},
			AmbientMap:  o.Material.AmbientMap,
			DiffuseMap:  o.Material.DiffuseMap,
			SpecularMap: o.Material.SpecularMap,
		})
		if err != nil {
			return err
		}
		data = append(data, &MeshData{
			Material:  mat,
			Vertices:  o.Vertices,
			Normals:   o.Normals,
			TexCoords: o.TexCoords,
		})
	}

	m.LoadFromData(data...)

	return nil
}

// LoadFromData loads a mesh from an array of MeshData
func (m *Mesh) LoadFromData(data ...*MeshData) error {
	const F = C.sizeof_float

	for _, d := range data {
		g := &meshGroup{
			material: d.Material,
		}

		g.count = int32(len(d.Vertices))
		hasNorms := len(d.Normals) > 0
		hasTxcds := len(d.TexCoords) > 0

		buf := make([]float32, 0, (len(d.Vertices)*3)+(len(d.Normals)*3)+(len(d.TexCoords)*2))
		for i := range d.Vertices {
			buf = append(buf, d.Vertices[i][0], d.Vertices[i][1], d.Vertices[i][2])
			if hasNorms {
				buf = append(buf, d.Normals[i][0], d.Normals[i][1], d.Normals[i][2])
			}
			if hasTxcds {
				buf = append(buf, d.TexCoords[i][0], d.TexCoords[i][1])
			}
		}

		g.size = len(buf)

		stride := int32(3 * F)
		if hasNorms {
			stride += int32(3 * F)
		}
		if hasTxcds {
			stride += int32(2 * F)
		}

		offset := 0

		gl.GenVertexArrays(1, &g.vao)
		gl.BindVertexArray(g.vao)

		gl.GenBuffers(1, &g.vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, g.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(buf)*F, gl.Ptr(buf), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(PositionAttrID)
		gl.VertexAttribPointer(PositionAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
		offset += 3 * F

		if hasNorms {
			gl.EnableVertexAttribArray(NormalAttrID)
			gl.VertexAttribPointer(NormalAttrID, 3, gl.FLOAT, false, stride, gl.PtrOffset(offset))
			offset += 3 * F
		}

		if hasTxcds {
			gl.EnableVertexAttribArray(TexCoordAttrID)
			gl.VertexAttribPointer(TexCoordAttrID, 2, gl.FLOAT, false, stride, gl.PtrOffset(offset))
		}

		m.groups = append(m.groups, g)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return nil
}

func (m *Mesh) GetMaterials() []*Material {
	mats := []*Material{}
	for _, g := range m.groups {
		mats = append(mats, g.material)
	}
	return mats
}

// UpdateData sets the data in the existing buffer
func (m *Mesh) UpdateData(data ...*MeshData) error {
	const F = C.sizeof_float

	for i := 0; i < len(data); i++ {
		d := data[i]
		g := m.groups[i]

		g.count = int32(len(d.Vertices))
		hasNorms := len(d.Normals) > 0
		hasTxcds := len(d.TexCoords) > 0

		buf := make([]float32, 0, (len(d.Vertices)*3)+(len(d.Normals)*3)+(len(d.TexCoords)*2))
		for i := range d.Vertices {
			buf = append(buf, d.Vertices[i][0], d.Vertices[i][1], d.Vertices[i][2])
			if hasNorms {
				buf = append(buf, d.Normals[i][0], d.Normals[i][1], d.Normals[i][2])
			}
			if hasTxcds {
				buf = append(buf, d.TexCoords[i][0], d.TexCoords[i][1])
			}
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, g.vbo)

		if g.size == len(buf) {
			gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(buf)*F, gl.Ptr(buf))
		} else {
			g.size = len(buf)
			gl.BufferData(gl.ARRAY_BUFFER, len(buf)*F, gl.Ptr(buf), gl.STATIC_DRAW)
		}

		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	return nil
}

// Render renders a Mesh to the screen
func (m *Mesh) Render(s *Shader) {
	for _, g := range m.groups {
		if g.material != nil {
			g.material.Bind(s)
		}

		gl.BindVertexArray(g.vao)
		gl.DrawArrays(gl.TRIANGLES, 0, g.count)

		if g.material != nil {
			g.material.UnBind()
		}
	}
}
