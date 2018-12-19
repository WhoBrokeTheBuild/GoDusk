package dusk

// IEntity is an Entity interface
type IEntity interface {
	Init(ILayer)
	Delete()

	GetLayer() ILayer

	AddComponent(IComponent)
	RemoveComponent(IComponent)

	Transform() *Transform
	SetTransform(*Transform)

	Update(*UpdateContext)
	Render(*RenderContext)
}

// Entity is an object with a Transform
type Entity struct {
	layer      ILayer
	transform  *Transform
	components []IComponent
}

// NewEntity returns a new, initialized Entity
func NewEntity(layer ILayer) *Entity {
	e := &Entity{}
	e.Init(layer)
	return e
}

// Init resets all data for the Entity
func (e *Entity) Init(layer ILayer) {
	e.Delete()

	e.layer = layer
	e.transform = NewTransform()
	e.components = []IComponent{}
}

// Delete frees all resources owned by the Entity
func (e *Entity) Delete() {
	e.layer = nil
	e.transform = nil
}

func (e *Entity) GetLayer() ILayer {
	return e.layer
}

func (e *Entity) AddComponent(component IComponent) {
	e.components = append(e.components, component)
}

func (e *Entity) RemoveComponent(component IComponent) {
	for i := 0; i < len(e.components); i++ {
		if e.components[i] == component {
			e.components = append(e.components[:i], e.components[i+1:]...)
		}
	}
}

// Transform returns the current transform
func (e *Entity) Transform() *Transform {
	return e.transform
}

// SetTransform sets the current transform
func (e *Entity) SetTransform(t *Transform) {
	e.transform = t
}

// Update fulfills the IEntity interface
func (e *Entity) Update(ctx *UpdateContext) {
	for _, c := range e.components {
		c.Update(ctx)
	}
}

// Render renders a Entity to the screen
func (e *Entity) Render(ctx *RenderContext) {
	for _, c := range e.components {
		c.Render(ctx)
	}
}
