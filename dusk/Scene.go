package dusk

// IScene is a Scene interface
type IScene interface {
	Init()
	Delete()

	AddActor(IActor)

	Update(*UpdateContext)
	Render(*RenderContext)
}

// Scene is a basic Scene
type Scene struct {
	actors []IActor
}

// NewScene returns a new, initialized Scene
func NewScene() *Scene {
	s := &Scene{}
	s.Init()
	return s
}

// Init initializes a Scene
func (s *Scene) Init() {
	s.Delete()

	s.actors = []IActor{}
}

// Delete clears all resources owned by the Scene
func (s *Scene) Delete() {
	s.actors = []IActor{}
}

// AddActor adds a new Actor to the scene
func (s *Scene) AddActor(actor IActor) {
	s.actors = append(s.actors, actor)
}

// Update calls Update() on all Actors
func (s *Scene) Update(ctx *UpdateContext) {
	for i := 0; i < len(s.actors); i++ {
		a := s.actors[i]
		a.Update(ctx)

		rb := a.RigidBody()
		if rb != nil {
			for j := i + 1; j < len(s.actors); j++ {
				if s.actors[j].RigidBody() != nil {
					a.RigidBody().CheckCollide(s.actors[j].RigidBody())
				}
			}
		}
	}
}

// Render calls Render() on all Actors
func (s *Scene) Render(ctx *RenderContext) {
	for _, a := range s.actors {
		a.Render(ctx)
	}
}
