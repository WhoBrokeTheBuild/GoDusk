package dusk

import (
	"fmt"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// UILayer represents a UI layer
type UILayer struct {
	Layer

	Buffer    *Texture
	Shader    *UIShader
	Mesh      *Mesh
	RenderCtx RenderContext
	Size      Vec2i

	frameID           uint32
	depthID           uint32
	needTextureUpdate bool
}

// NewUILayer returns a new UILayer of the given size
func NewUILayer(app *App) (*UILayer, error) {
	size := Vec2i{app.Window.Width, app.Window.Height}

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

	return &UILayer{
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
func (ui *UILayer) Delete() {
	if ui.Buffer != nil {
		ui.Buffer.Delete()
		ui.Buffer = nil
	}

	if ui.Shader != nil {
		ui.Shader.Delete()
		ui.Shader = nil
	}

	if ui.Mesh != nil {
		ui.Mesh.Delete()
		ui.Mesh = nil
	}

	if ui.depthID != 0 {
		gl.DeleteRenderbuffers(1, &ui.depthID)
		ui.depthID = 0
	}

	if ui.frameID != 0 {
		gl.DeleteFramebuffers(1, &ui.frameID)
		ui.frameID = 0
	}
}

// Render renders the current buffer to the screen
func (ui *UILayer) Render(_ *RenderContext) {
	ui.Shader.Bind(&ui.RenderCtx, nil)
	gl.UniformMatrix4fv(ui.Shader.UniformLocation("uProjection"), 1, false, &ui.RenderCtx.Projection[0])

	gl.BindFramebuffer(gl.FRAMEBUFFER, ui.frameID)
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, e := range ui.GetEntities() {
		e.Render(&ui.RenderCtx)
		gl.Clear(gl.DEPTH_BUFFER_BIT)
	}

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)

	gl.Uniform1i(ui.Shader.UniformLocation("uTexture"), 0)
	gl.ActiveTexture(gl.TEXTURE0)
	ui.Buffer.Bind()

	gl.Clear(gl.DEPTH_BUFFER_BIT)
	ui.Mesh.Render(ui.Shader)

	gl.BindTexture(gl.TEXTURE_2D, 0)
}
