package dusk

import (
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/asset"
	"github.com/WhoBrokeTheBuild/GoDusk/context"
	"github.com/WhoBrokeTheBuild/GoDusk/types"
	"github.com/WhoBrokeTheBuild/GoDusk/ui"
)

// AppOptions is used to create a new App
type AppOptions struct {
	Window *WindowOptions
}

// DefaultAppOptions returns the default values for AppOptions
func DefaultAppOptions() *AppOptions {
	return &AppOptions{
		Window: DefaultWindowOptions(),
	}
}

// UpdateFunc is a function meant to be called during Update
type UpdateFunc func(*context.Update)

// RenderFunc is a function meant to be called during Render
type RenderFunc func(*context.Render)

// App represents an application
type App struct {
	Window *Window
	UI     *ui.Overlay

	defaultShader *asset.Shader

	updateFuncs []UpdateFunc
	renderFuncs []RenderFunc

	updateCtx *context.Update
	renderCtx *context.Render
}

// NewApp creates an new App from the given AppOptions
func NewApp(opts *AppOptions) (app *App, err error) {
	app = &App{}

	app.Window, err = NewWindow(opts.Window)
	if err != nil {
		app.Delete()
		return
	}

	app.Window.RegisterResizeFunc(func(width, height int) {})

	app.UI, err = ui.NewOverlay(types.Vec2i{app.Window.Width, app.Window.Height})
	if err != nil {
		app.Delete()
		return
	}

	app.defaultShader, err = asset.NewShaderFromFiles([]string{"data/shaders/default.vs.glsl", "data/shaders/default.fs.glsl"})
	if err != nil {
		return
	}

	aspect := float32(app.Window.Width) / float32(app.Window.Height)
	app.updateCtx = &context.Update{}
	app.renderCtx = &context.Render{
		View:       mgl32.LookAtV(mgl32.Vec3{2, 2, 2}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}),
		Projection: mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0),
		Shader:     app.defaultShader,
	}

	return
}

// Delete frees an App's resources
func (app *App) Delete() {
	if app.defaultShader != nil {
		app.defaultShader.Delete()
		app.defaultShader = nil
	}

	if app.UI != nil {
		app.UI.Delete()
		app.UI = nil
	}

	if app.Window != nil {
		app.Window.Delete()
		app.Window = nil
	}
}

func (app *App) RegisterUpdateFunc(fun UpdateFunc) {
	app.updateFuncs = append(app.updateFuncs, fun)
}

func (app *App) RegisterRenderFunc(fun RenderFunc) {
	app.renderFuncs = append(app.renderFuncs, fun)
}

// Run starts the update/render loop for the App, it will not return until the window closes
func (app *App) Run() {
	const (
		frameDelay = 1.0 / 60.0
		fpsDelay   = 1.0
	)

	var (
		frameElap  = 0.0
		frameCount = 0
		fpsElap    = 0.0
	)

	prev := glfw.GetTime()
	for !app.Window.ShouldClose() {
		time := glfw.GetTime()
		elapsed := time - prev
		prev = time

		frameElap += elapsed
		fpsElap += elapsed
		app.updateCtx.DeltaTime = 0.0
		app.updateCtx.ElapsedTime = elapsed

		if fpsElap >= fpsDelay {
			app.updateCtx.FPS = frameCount

			fpsElap = 0.0
			frameCount = 0
		}

		app.Window.PollEvents()

		for _, f := range app.updateFuncs {
			f(app.updateCtx)
		}

		if frameElap >= frameDelay {
			frameCount++
			frameElap = 0.0

			gl.ClearColor(0.0, 0.4, 0.8, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			for _, f := range app.renderFuncs {
				f(app.renderCtx)
			}

			app.UI.Draw()

			app.Window.SwapBuffers()
		}
	}
}
