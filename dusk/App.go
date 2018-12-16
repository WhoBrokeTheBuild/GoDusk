package dusk

import (
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
type UpdateFunc func(*UpdateContext)

// RenderFunc is a function meant to be called during Render
type RenderFunc func(*RenderContext)

// App represents an application
type App struct {
	Window *Window
	UI     *UI
	Scene  *Scene

	defaultCamera *Camera

	updateFuncs []UpdateFunc
	renderFuncs []RenderFunc

	updateCtx *UpdateContext
	renderCtx *RenderContext
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

	app.Scene = NewScene()

	app.UI, err = NewUI(Vec2i{app.Window.Width, app.Window.Height})
	if err != nil {
		app.Delete()
		return
	}

	app.defaultCamera = NewCamera(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0})

	aspect := float32(app.Window.Width) / float32(app.Window.Height)
	app.updateCtx = &UpdateContext{}
	app.renderCtx = &RenderContext{
		Projection: mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 10000.0),
		Camera:     app.defaultCamera,
	}

	return
}

// Delete frees an App's resources
func (app *App) Delete() {
	if app.UI != nil {
		app.UI.Delete()
		app.UI = nil
	}

	if app.Window != nil {
		app.Window.Delete()
		app.Window = nil
	}
}

// RegisterUpdateFunc adds an update callback function
func (app *App) RegisterUpdateFunc(fun UpdateFunc) {
	app.updateFuncs = append(app.updateFuncs, fun)
}

// RegisterRenderFunc adds a render callback function
func (app *App) RegisterRenderFunc(fun RenderFunc) {
	app.renderFuncs = append(app.renderFuncs, fun)
}

// SetScene sets the current Scene
func (app *App) SetScene(scene *Scene) {
	app.Scene = scene
}

// GetRenderContext returns the default Render Context
func (app *App) GetRenderContext() *RenderContext {
	return app.renderCtx
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
		app.updateCtx.DeltaTime = float32(elapsed / frameDelay)
		app.updateCtx.ElapsedTime = elapsed

		if fpsElap >= fpsDelay {
			app.updateCtx.FPS = frameCount

			fpsElap = 0.0
			frameCount = 0
		}

		app.Window.PollEvents()

		if app.Scene != nil {
			app.Scene.Update(app.updateCtx)
		}

		for _, f := range app.updateFuncs {
			f(app.updateCtx)
		}

		if frameElap >= frameDelay {
			frameCount++
			frameElap = 0.0

			gl.ClearColor(0.0, 0.4, 0.8, 1.0)
			gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

			if app.Scene != nil {
				app.Scene.Render(app.renderCtx)
			}

			for _, f := range app.renderFuncs {
				f(app.renderCtx)
			}

			app.UI.Draw()

			app.Window.SwapBuffers()
		}
	}
}
