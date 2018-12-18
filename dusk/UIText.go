package dusk

import (
	"image"
	"image/color"
	"path/filepath"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

var _fonts map[string]*truetype.Font

func init() {
	_fonts = map[string]*truetype.Font{}
}

// UIText is a UIElement that draws text to the screen
type UIText struct {
	UIImage

	Text  string
	Size  float64
	Color color.Color
	Font  *truetype.Font
	Face  font.Face
}

// NewUIText returns a new UIText from a given string, font, font size, and color
func NewUIText(text string, font string, size float64, color color.Color) *UIText {
	var f *truetype.Font

	font = filepath.Clean(font)
	if tmp, found := _fonts[font]; found {
		Loadf("ui.Font [%v]+", font)
		f = tmp
	} else {
		Loadf("ui.Font [%v]", font)
		b, err := Load(font)
		if err != nil {
			Errorf("%v", err)
			return nil
		}

		f, err = freetype.ParseFont(b)
		if err != nil {
			Errorf("%v", err)
			return nil
		}
		_fonts[font] = f
	}

	c := &UIText{
		Text:  text,
		Size:  size,
		Color: color,
		Font:  f,
	}
	c.updateUITexture()

	return c
}

// SetText sets the text to be rendered
func (c *UIText) SetText(text string) {
	c.Text = text
	c.updateUITexture()
}

func (c *UIText) updateUITexture() {
	var err error

	if c.Texture != nil {
		c.Texture.Delete()
	}

	d := &font.Drawer{
		Dst: nil,
		Src: image.NewUniform(c.Color),
		Dot: freetype.Pt(0, int(c.Size)),
		Face: truetype.NewFace(c.Font, &truetype.Options{
			Size:    c.Size,
			DPI:     72,
			Hinting: font.HintingFull,
		}),
	}

	m := d.MeasureString(c.Text)
	buffer := image.NewRGBA(image.Rect(0, 0, m.Ceil(), int(c.Size*1.5)))

	d.Dst = buffer
	d.DrawString(c.Text)

	s := buffer.Rect.Size()
	c.Texture, err = NewTextureFromData(buffer.Pix, gl.RGBA, gl.RGBA, s.X, s.Y)
	if err != nil {
		Errorf("%v", err)
	}

	c.SetSize(mgl32.Vec2{float32(s.X), float32(s.Y)})
}
