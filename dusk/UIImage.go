package dusk

import (

	// JPEG support
	_ "image/jpeg"

	// PNG support
	_ "image/png"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// UIImage is a UIElement that draws an image to the screen
type UIImage struct {
	BaseUIElement
	Bounds  mgl32.Vec4
	Texture *Texture
	Mesh    *Mesh
}

// NewUIImageFromFile returns a new UIImage from the given file
func NewUIImageFromFile(filename string) *UIImage {
	c := &UIImage{}
	err := c.LoadFromFile(filename)
	if err != nil {
		c.Delete()
		Errorf("%v", err)
		return nil
	}
	return c
}

// NewUIImageFromData returns a new UIImage from the given data, width, and height
func NewUIImageFromData(data []uint8, intFormat uint32, format int32, width, height int) *UIImage {
	c := &UIImage{}
	err := c.LoadFromData(data, intFormat, format, width, height)
	if err != nil {
		c.Delete()
		Errorf("%v", err)
		return nil
	}
	return c
}

// Delete frees all resources owned by the UIImage
func (c *UIImage) Delete() {
	if c.Texture != nil {
		c.Texture.Delete()
		c.Texture = nil
	}
}

// LoadFromFile loads an UIImage from the given file
func (c *UIImage) LoadFromFile(filename string) error {
	var err error
	c.Delete()

	c.Texture, err = NewTextureFromFile(filename)
	if err != nil {
		c.Delete()
		return err
	}

	c.SetSize(c.Texture.Size)

	return nil
}

// LoadFromData loads an UIImage from the given data, width, and height
func (c *UIImage) LoadFromData(data []uint8, intFormat uint32, format int32, width, height int) error {
	var err error
	c.Delete()

	c.Texture, err = NewTextureFromData(data, intFormat, format, width, height)
	if err != nil {
		c.Delete()
		return err
	}

	c.SetSize(c.Texture.Size)

	return nil
}

// SetPosition sets the UIImage's position
func (c *UIImage) SetPosition(pos mgl32.Vec2) {
	c.BaseUIElement.SetPosition(pos)
	c.updateMesh()
}

// SetSize sets the UIImage's size
func (c *UIImage) SetSize(size mgl32.Vec2) {
	c.BaseUIElement.SetSize(size)
	c.updateMesh()
}

func (c *UIImage) updateMesh() {
	var err error
	pos := c.GetPosition()
	size := c.GetSize()

	x := pos.X()
	y := pos.Y()
	w := size.X()
	h := size.Y()

	if c.Mesh == nil {
		c.Mesh, err = new2DMesh(
			mgl32.Vec4{x, y, x + w, y + h},
			mgl32.Vec4{0, 0, 1, 1})
		if err != nil {
			c.Delete()
			Errorf("%v", err)
		}
	} else {
		err = update2DMesh(c.Mesh,
			mgl32.Vec4{x, y, x + w, y + h},
			mgl32.Vec4{0, 0, 1, 1})
		if err != nil {
			c.Delete()
			Errorf("%v", err)
		}
	}
}

// Draw renders the UIImage to the buffer
func (c *UIImage) Draw(ctx *RenderContext) {
	s := GetUIShader()

	gl.Uniform1i(s.UniformLocation("uTexture"), 0)

	gl.ActiveTexture(gl.TEXTURE0)
	if c.Texture != nil {
		c.Texture.Bind()
	}
	if c.Mesh != nil {
		c.Mesh.Render(s)
	}
}
