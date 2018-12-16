package dusk

// Actor is an object with a Transform
type Actor struct {
	Transform *Transform
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

// Update updates an Actor
func (a *Actor) Update(ctx *UpdateContext) {

}

// Render renders a Actor to the screen
func (a *Actor) Render(ctx *RenderContext) {
	s := GetDefaultShader()
	s.Bind(ctx, a.Transform.GetMatrix())

	for _, m := range a.Meshes {
		m.Render(s)
	}
}

// DefaultShader is the default shader used to render meshes
type DefaultShader struct {
	Shader
}

var _defaultShader *DefaultShader
