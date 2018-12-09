package dusk

import (
	"runtime"

	// Includes support for .png icons
	_ "image/png"
	// Includes support for .jpg, .jpeg icons
	_ "image/jpeg"
	// Includes support for .gif icons
	_ "image/gif"

	//"github.com/go-gl/mathgl/mgl32"

	"github.com/WhoBrokeTheBuild/GoDusk/log"
)

const (
	// Version is the GoDusk Version
	Version = "0.1.0"

	// GLMajor is the OpenGL Major Version
	GLMajor = 4
	// GLMinor is the OpenGL Minor Version
	GLMinor = 1
)

func init() {
	runtime.LockOSThread()

	runtime.GOMAXPROCS(runtime.NumCPU())
	log.Infof("CPU Cores: %d", runtime.NumCPU())
}
