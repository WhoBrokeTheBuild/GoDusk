package dusk

// IActor is an Actor interface
type IActor interface {
	Init()
	Delete()

	Transform() *Transform
	SetTransform(*Transform)

	Shader() IShader
	SetShader(IShader)

	Update(*UpdateContext)
	Render(*RenderContext)
}

// Actor is an object with a Transform
type Actor struct {
	transform *Transform
	meshes    []*Mesh
	shader    IShader
}

// NewActor returns a new, initialized Actor
func NewActor() *Actor {
	a := &Actor{}
	a.Init()
	return a
}

// Init resets all data for the Actor
func (a *Actor) Init() {
	a.Delete()

	a.transform = NewTransform()
	a.meshes = []*Mesh{}
	a.shader = GetDefaultShader()
}

// Transform returns the current transform
func (a *Actor) Transform() *Transform {
	return a.transform
}

// SetTransform sets the current transform
func (a *Actor) SetTransform(t *Transform) {
	a.transform = t
}

// Shader returns the current shader
func (a *Actor) Shader() IShader {
	return a.shader
}

// SetShader sets the current shader
func (a *Actor) SetShader(s IShader) {
	a.shader = s
}

// AddMesh adds one or more meshes to the Actor
func (a *Actor) AddMesh(mesh ...*Mesh) {
	a.meshes = append(a.meshes, mesh...)
}

// Delete frees all resources owned by the Actor
func (a *Actor) Delete() {
	a.meshes = []*Mesh{}
	a.shader = nil
}

// Update fulfills the IActor interface
func (a *Actor) Update(ctx *UpdateContext) {}

// Render renders a Actor to the screen
func (a *Actor) Render(ctx *RenderContext) {
	s := a.Shader()
	s.Bind(ctx, a.transform.GetMatrix())

	for _, m := range a.meshes {
		m.Render(s)
	}
}
