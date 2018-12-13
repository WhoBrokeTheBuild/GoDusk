package dusk

type Scene struct {
	Actors []*Actor
}

func NewScene() *Scene {
	return &Scene{}
}

func (s *Scene) AddActor(actor *Actor) {
	s.Actors = append(s.Actors, actor)
}

func (s *Scene) Update(ctx *UpdateContext) {
	for _, a := range s.Actors {
		a.Update(ctx)
	}
}

func (s *Scene) Render(ctx *RenderContext) {
	for _, a := range s.Actors {
		a.Render(ctx)
	}
}
