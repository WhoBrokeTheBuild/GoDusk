package dusk

// ILayer is a Layer interface
type ILayer interface {
	Init()
	Delete()

	AddEntity(IEntity)
	RemoveEntity(IEntity)

	Update(*UpdateContext)
	Render(*RenderContext)

	GetEntities() []IEntity
}

// Layer is a basic Layer
type Layer struct {
	entities []IEntity
}

// NewLayer returns a new, initialized Layer
func NewLayer() *Layer {
	s := &Layer{}
	s.Init()
	return s
}

// Init initializes a Layer
func (s *Layer) Init() {
	s.Delete()

	s.entities = []IEntity{}
}

// Delete clears all resources owned by the Layer
func (s *Layer) Delete() {
	s.entities = []IEntity{}
}

// AddEntity adds a new Entity to the Layer
func (s *Layer) AddEntity(entity IEntity) {
	s.entities = append(s.entities, entity)
}

// RemoveEntity adds a new Entity to the Layer
func (s *Layer) RemoveEntity(entity IEntity) {
	for i := 0; i < len(s.entities); i++ {
		if s.entities[i] == entity {
			s.entities = append(s.entities[:i], s.entities[i+1:]...)
		}
	}
}

// Update calls Update() on all entities
func (s *Layer) Update(ctx *UpdateContext) {
	for _, e := range s.entities {
		e.Update(ctx)
	}
}

// Render calls Render() on all entities
func (s *Layer) Render(ctx *RenderContext) {
	for _, e := range s.entities {
		e.Render(ctx)
	}
}

func (s *Layer) GetEntities() []IEntity {
	return s.entities
}
