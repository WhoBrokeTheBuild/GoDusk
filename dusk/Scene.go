package dusk

import "github.com/WhoBrokeTheBuild/GoDusk/context"

type Scene interface {
	Update(*context.Update)
	Render(*context.Render)
}
