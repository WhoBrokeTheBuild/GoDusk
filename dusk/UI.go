package dusk

import (
	"fmt"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// UI represents a UI layer
type UI struct {
	Buffer     *Texture
	Shader     *UIShader
	Mesh       *Mesh
	RenderCtx  RenderContext
	Size       Vec2i
	UIElements []UIElement

	frameID           uint32
	depthID           uint32
	needTextureUpdate bool
}

// NewUI returns a new UI of the given size
func NewUI(size Vec2i) (*UI, error) {
	mesh, err := new2DMesh(mgl32.Vec4{0, 0, float32(size.X()), float32(size.Y())}, mgl32.Vec4{0, 1, 1, 0})
	if err != nil {
		return nil, err
	}

	var (
		frameID uint32
		depthID uint32
	)

	buffer, err := NewTextureFromData(nil, gl.RGBA8, gl.RGBA, int(size.X()), int(size.Y()))
	if err != nil {
		return nil, err
	}

	gl.GenRenderbuffers(1, &depthID)
	gl.BindRenderbuffer(gl.RENDERBUFFER, depthID)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH_COMPONENT, int32(size.X()), int32(size.Y()))
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)

	gl.GenFramebuffers(1, &frameID)
	gl.BindFramebuffer(gl.FRAMEBUFFER, frameID)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, buffer.ID, 0)
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_ATTACHMENT, gl.RENDERBUFFER, depthID)

	status := gl.CheckFramebufferStatus(gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		return nil, fmt.Errorf("Failed to create Framebuffer")
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	return &UI{
		Buffer: buffer,
		Shader: GetUIShader(),
		Mesh:   mesh,
		Size:   size,

		RenderCtx: RenderContext{
			Projection: mgl32.Ortho2D(0, float32(size.X()), 0, float32(size.Y())),
		},

		frameID: frameID,
		depthID: depthID,

		needTextureUpdate: true,
	}, nil
}

// Delete frees all resources owned by the UI
func (o *UI) Delete() {
	if o.Buffer != nil {
		o.Buffer.Delete()
		o.Buffer = nil
	}

	if o.Shader != nil {
		o.Shader.Delete()
		o.Shader = nil
	}

	if o.Mesh != nil {
		o.Mesh.Delete()
		o.Mesh = nil
	}

	if o.depthID != 0 {
		gl.DeleteRenderbuffers(1, &o.depthID)
		o.depthID = 0
	}

	if o.frameID != 0 {
		gl.DeleteFramebuffers(1, &o.frameID)
		o.frameID = 0
	}
}

// Update processes all events
func (o *UI) Update(ctx *UpdateContext) {

}

// Draw renders the current buffer to the screen
func (o *UI) Draw() {
	o.Shader.Bind(&o.RenderCtx, nil)
	gl.UniformMatrix4fv(o.Shader.UniformLocation("uProjection"), 1, false, &o.RenderCtx.Projection[0])

	gl.BindFramebuffer(gl.FRAMEBUFFER, o.frameID)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	for _, c := range o.UIElements {
		c.Draw(&o.RenderCtx)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
	}
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	gl.Uniform1i(o.Shader.UniformLocation("uTexture"), 0)
	gl.ActiveTexture(gl.TEXTURE0)
	o.Buffer.Bind()

	gl.Clear(gl.DEPTH_BUFFER_BIT)
	o.Mesh.Render(o.Shader)

	gl.BindTexture(gl.TEXTURE_2D, 0)
}

// AddElement adds the given UIElement as a child node of the UI
func (o *UI) AddElement(c UIElement) {
	o.UIElements = append(o.UIElements, c)
}
