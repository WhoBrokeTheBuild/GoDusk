package dusk

import (
	"bytes"
	"image"

	// Includes support for .png icons
	_ "image/png"
	// Includes support for .jpg, .jpeg icons
	_ "image/jpeg"
	// Includes support for .gif icons
	_ "image/gif"

	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/WhoBrokeTheBuild/GoDusk/load"
	"github.com/WhoBrokeTheBuild/GoDusk/log"
)

// WindowOptions is used to create a new Window
type WindowOptions struct {
	Width  int
	Height int
	Title  string
	Icons  []string
}

// DefaultWindowOptions returns the default values for WindowOptions
func DefaultWindowOptions() *WindowOptions {
	return &WindowOptions{
		Width:  1024,
		Height: 768,
		Title:  "GoDusk",
		Icons:  []string{},
	}
}

type ResizeFunc func(int, int)

// Window represents a Window
type Window struct {
	Width  int
	Height int
	Title  string

	resizeFuncs []ResizeFunc

	glfwWindow *glfw.Window
}

// NewWindow creates a new Window from the given WindowOptions
func NewWindow(opts *WindowOptions) (w *Window, err error) {
	w = &Window{
		Width:  opts.Width,
		Height: opts.Height,
		Title:  opts.Title,
	}

	err = glfw.Init()
	if err != nil {
		w.Delete()
		return
	}

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, GLMajor)
	glfw.WindowHint(glfw.ContextVersionMinor, GLMinor)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	w.glfwWindow, err = glfw.CreateWindow(opts.Width, opts.Height, opts.Title, nil, nil)
	if err != nil {
		w.Delete()
		return
	}

	w.glfwWindow.SetSizeCallback(func(_ *glfw.Window, width, height int) {
		w.Width = width
		w.Height = height

		width, height = w.glfwWindow.GetFramebufferSize()
		gl.Viewport(0, 0, int32(width), int32(height))

		for _, f := range w.resizeFuncs {
			f(width, height)
		}
	})

	if len(opts.Icons) > 0 {
		icons := []image.Image{}
		for _, file := range opts.Icons {
			log.Loadf("Icon [%v]", file)
			b, err := load.Load(file)
			if err != nil {
				log.Warnf("Failed to find icon [%v]", file)
				continue
			}
			image, _, err := image.Decode(bytes.NewReader(b))
			if err != nil {
				log.Warnf("Failed to decode icon [%v]", file)
				continue
			}
			icons = append(icons, image)
		}
		w.glfwWindow.SetIcon(icons)
	}

	w.glfwWindow.MakeContextCurrent()
	glfw.SwapInterval(1)

	err = gl.Init()
	if err != nil {
		w.Delete()
		return
	}

	log.Infof("OpenGL Version: [%s]", gl.GoStr(gl.GetString(gl.VERSION)))
	log.Infof("GLSL Version: [%s]", gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION)))
	log.Infof("OpenGL Vendor: [%s]", gl.GoStr(gl.GetString(gl.VENDOR)))
	log.Infof("OpenGL Renderer: [%s]", gl.GoStr(gl.GetString(gl.RENDERER)))

	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	w.glfwWindow.SwapBuffers()

	return
}

// Delete frees a Window's resources
func (w *Window) Delete() {
	if w.glfwWindow != nil {
		w.glfwWindow.Destroy()
		w.glfwWindow = nil
	}

	glfw.Terminate()
}

func (w *Window) RegisterResizeFunc(fun ResizeFunc) {
	w.resizeFuncs = append(w.resizeFuncs, fun)
}

// ShouldClose returns true if the current window should close
func (w *Window) ShouldClose() bool {
	if w.glfwWindow == nil {
		return true
	}
	return w.glfwWindow.ShouldClose()
}

// PollEvents checks for window events and dispatches handlers
func (w *Window) PollEvents() {
	glfw.PollEvents()
}

// SwapBuffers swaps the back and front buffers
func (w *Window) SwapBuffers() {
	w.glfwWindow.SwapBuffers()
}

// SetTitle sets the window title
func (w *Window) SetTitle(title string) {
	w.Title = title
	w.glfwWindow.SetTitle(w.Title)
}

// SetSize sets the window size
func (w *Window) SetSize(width, height int) {
	w.glfwWindow.SetSize(width, height)
}
