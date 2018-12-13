package dusk

import (
	"github.com/go-gl/mathgl/mgl32"
)

// UIElement is the interface type for all UI UIElements
type UIElement interface {
	Draw(*RenderContext)
	Delete()
}

// BaseUIElement is a stub UIElement made to be inherited from
type BaseUIElement struct {
	Position mgl32.Vec2
	Size     mgl32.Vec2
}

// Delete frees all resources owned by the UIElement
func (c *BaseUIElement) Delete() {
}

// GetPosition returns the UIElement's current position
func (c *BaseUIElement) GetPosition() mgl32.Vec2 {
	return c.Position
}

// SetPosition sets the UIElement's position
func (c *BaseUIElement) SetPosition(pos mgl32.Vec2) {
	c.Position = pos
}

// GetSize returns the UIElement's current size
func (c *BaseUIElement) GetSize() mgl32.Vec2 {
	return c.Size
}

// SetSize sets the UIElement's size
func (c *BaseUIElement) SetSize(size mgl32.Vec2) {
	c.Size = size
}

// Draw renders a BaseUIElement
func (c *BaseUIElement) Draw(ctx *RenderContext) {
}
