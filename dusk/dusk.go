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
)

import "io/ioutil"

//go:generate go-bindata -tags debug -debug -pkg $GOPACKAGE -o data-debug.gen.go data/...
//go:generate go-bindata -tags !debug -pkg $GOPACKAGE -o data.gen.go data/...

const (
	// Version is the GoDusk Version
	Version = "0.1.0"

	// GLMajor is the OpenGL Major Version
	GLMajor = 4
	// GLMinor is the OpenGL Minor Version
	GLMinor = 1

	// InvalidID is an invalid OpenGL ID
	InvalidID uint32 = 0
)

func init() {
	runtime.LockOSThread()

	runtime.GOMAXPROCS(runtime.NumCPU())
	Infof("CPU Cores: %d", runtime.NumCPU())
}

// LoadFunc represents a function to load an asset, e.g. ioutil.ReadFile
type LoadFunc func(string) ([]byte, error)

// Functions to use to load assets, tried in reverse order
var _loadFuncs = []LoadFunc{
	ioutil.ReadFile,
	Asset,
}

// Load attempts to load an asset from all registered asset funcs
func Load(filename string) (b []byte, err error) {
	for i := range _loadFuncs {
		load := _loadFuncs[len(_loadFuncs)-1-i]
		b, err = load(filename)
		if err == nil {
			return
		}
	}
	return
}

// RegisterFunc prepends a new asset loading function to the list of functions to try
func RegisterFunc(fun LoadFunc) {
	_loadFuncs = append(_loadFuncs, fun)
}
