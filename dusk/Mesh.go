package dusk

import (
	"C"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// Mesh represents a set of OpenGL Vertex Array Objects and Material data
type Mesh struct {
	name     string
	material *Material
	vao      uint32
	vbo      uint32
	size     int
	count    int32
}

// MeshData is the intermediate data format for loading Meshes from Memory
type MeshData struct {
	Name     string
	Material *Material

	Vertices  []mgl32.Vec3
	Normals   []mgl32.Vec3
	TexCoords []mgl32.Vec2
}

// NewMeshFromData returns a new Mesh from the given MeshData
func NewMeshFromData(data *MeshData) (*Mesh, error) {
	m := &Mesh{}

	err := m.LoadFromData(data)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, err
}

// Delete frees all resources owned by the Mesh
func (m *Mesh) Delete() {
	if m.material != nil {
		m.material.Delete()
	}
	if m.vbo != InvalidID {
		gl.DeleteBuffers(1, &m.vbo)
		m.vbo = InvalidID
	}
	if m.vao != InvalidID {
		gl.DeleteVertexArrays(1, &m.vao)
		m.vao = InvalidID
	}
}

// LoadFromData loads a mesh from an array of MeshData
func (m *Mesh) LoadFromData(data *MeshData) error {
	const F = C.sizeof_float

	m.material = data.Material

	m.count = int32(len(data.Vertices))
	hasNorms := len(data.Normals) > 0
	hasTxcds := len(data.TexCoords) > 0

	buf := make([]float32, 0, (len(data.Vertices)*3)+(len(data.Normals)*3)+(len(data.TexCoords)*2))
	for i := range data.Vertices {
		buf = append(buf, data.Vertices[i][0], data.Vertices[i][1], data.Vertices[i][2])
		if hasNorms {
			buf = append(buf, data.Normals[i][0], data.Normals[i][1], data.Normals[i][2])
		}
		if hasTxcds {
			buf = append(buf, data.TexCoords[i][0], data.TexCoords[i][1])
		}
	}

	m.size = len(buf)

	stride := int32(3 * F)
	if hasNorms {
		stride += int32(3 * F)
	}
	if hasTxcds {
		stride += int32(2 * F)
	}

	offset := 0

	gl.GenVertexArrays(1, &m.vao)
	gl.BindVertexArray(m.vao)

	gl.GenBuffers(1, &m.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)
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

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return nil
}

func (m *Mesh) GetMaterial() *Material {
	return m.material
}

// UpdateData sets the data in the existing buffer
func (m *Mesh) UpdateData(data *MeshData) error {
	const F = C.sizeof_float

	if data.Material != nil {
		m.material = data.Material
	}

	m.count = int32(len(data.Vertices))
	hasNorms := len(data.Normals) > 0
	hasTxcds := len(data.TexCoords) > 0

	buf := make([]float32, 0, (len(data.Vertices)*3)+(len(data.Normals)*3)+(len(data.TexCoords)*2))
	for i := range data.Vertices {
		buf = append(buf, data.Vertices[i][0], data.Vertices[i][1], data.Vertices[i][2])
		if hasNorms {
			buf = append(buf, data.Normals[i][0], data.Normals[i][1], data.Normals[i][2])
		}
		if hasTxcds {
			buf = append(buf, data.TexCoords[i][0], data.TexCoords[i][1])
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, m.vbo)

	if m.size == len(buf) {
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(buf)*F, gl.Ptr(buf))
	} else {
		m.size = len(buf)
		gl.BufferData(gl.ARRAY_BUFFER, len(buf)*F, gl.Ptr(buf), gl.STATIC_DRAW)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	return nil
}

// Render renders a Mesh to the screen
func (m *Mesh) Render(s IShader) {
	if m.material != nil {
		m.material.Bind(s)
	}

	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.count)

	if m.material != nil {
		m.material.UnBind()
	}
}
