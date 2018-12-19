package dusk

import (
	"fmt"
	"path/filepath"
)

// ModelLoader is a function that loads mesh data
type ModelLoader func(filename string) ([]*MeshData, error)

type modelFormat struct {
	name   string
	exts   []string
	loader ModelLoader
}

func (f *modelFormat) hasExt(ext string) bool {
	for i := range f.exts {
		if f.exts[i] == ext {
			return true
		}
	}
	return false
}

var _modelFormats = map[string]modelFormat{}

// RegisterModelFormat adds a new handler for loading mesh files
func RegisterModelFormat(name string, exts []string, loader ModelLoader) {
	_modelFormats[name] = modelFormat{
		name:   name,
		exts:   exts,
		loader: loader,
	}
}

type Model struct {
	Component
	shader IShader
	meshes map[string]*Mesh
}

// NewModelFromFile returns a new Mesh from the given file
func NewModelFromFile(entity IEntity, filename string) (*Model, error) {
	m := &Model{
		shader: GetDefaultShader(),
		meshes: map[string]*Mesh{},
	}
	m.Init(entity)

	err := m.LoadFromFile(filename)
	if err != nil {
		m.Delete()
		return nil, err
	}

	return m, err
}

// LoadFromFile loads the meshes from a given file
func (m *Model) LoadFromFile(filename string) error {
	filename = filepath.Clean(filename)

	var loader ModelLoader

	ext := filepath.Ext(filename)
	for _, f := range _modelFormats {
		if f.hasExt(ext) {
			loader = f.loader
		}
	}

	if loader == nil {
		return fmt.Errorf("Unsupported format [%v]", ext)
	}

	Loadf("asset.Mesh [%v]", filename)
	data, err := loader(filename)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return fmt.Errorf("No data loaded from [%v]", filename)
	}

	for _, d := range data {
		m.meshes[d.Name], err = NewMeshFromData(d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Model) Render(ctx *RenderContext) {
	m.shader.Bind(ctx, m.GetEntity().Transform().GetMatrix())
	for _, mesh := range m.meshes {
		mesh.Render(m.shader)
	}
}
