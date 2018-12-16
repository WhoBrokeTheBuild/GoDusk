package dusk

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	gl "github.com/go-gl/gl/v4.1-core/gl"
)

const (
	// ShaderIncludePath is the path to search for files from #include <>
	ShaderIncludePath = "data/shaders/include"
)

var (
	_versionString  string
	_defaultDefines = map[string]string{}
)

// AddShaderDefines adds default #define values for processed shaders
func AddShaderDefines(defines map[string]interface{}) {
	for k, v := range defines {
		_defaultDefines[k] = fmt.Sprintf("%v", v)
	}
}

// GetShaderDefines returns a copy of the map of default shader defines
func GetShaderDefines() map[string]string {
	tmp := map[string]string{}
	for k, v := range _defaultDefines {
		tmp[k] = v
	}
	return tmp
}

// IShader is an interface representing an OpenGL Shader Program
type IShader interface {
	InitFromFiles(...string)
	InitFromData(...*ShaderData)
	Delete()

	ID() uint32

	Bind(*RenderContext, interface{})

	UniformLocation(string) int32
}

// Shader represents a generic shader
type Shader struct {
	id       uint32
	uniforms map[string]int32
}

// ShaderData represents a shader's code and type
type ShaderData struct {
	Code string
	Type uint32
}

// InitFromFiles loads a new shader from a set of filenames
func (s *Shader) InitFromFiles(filename ...string) {
	s.Delete()

	var err error
	s.id, err = loadShaderFromFiles(filename...)
	if err != nil {
		s.id = InvalidID
	}
}

// InitFromData loads a new shader from a set of filenames
func (s *Shader) InitFromData(data ...*ShaderData) {
	s.Delete()

	var err error
	s.id, err = loadShaderFromData(data...)
	if err != nil {
		s.id = InvalidID
	}
}

// Delete frees all resources owned by the Shader
func (s *Shader) Delete() {
	if s.id != InvalidID {
		gl.DeleteProgram(s.id)
		s.id = InvalidID
	}
}

// ID returns the underlying OpenGL Shader Program ID
func (s *Shader) ID() uint32 {
	return s.id
}

// Bind binds this shader and all uniforms
func (s *Shader) Bind(_ *RenderContext, _ interface{}) {
	gl.UseProgram(s.id)
}

// UniformLocation returns the location of the given uniform, or -1
func (s *Shader) UniformLocation(name string) int32 {
	if len(s.uniforms) == 0 {
		s.cacheUniforms()
	}
	if loc, found := s.uniforms[name]; found {
		return loc
	}
	return -1
}

func (s *Shader) cacheUniforms() {
	s.uniforms = map[string]int32{}

	var count int32
	var size int32
	var length int32
	var tp uint32

	buf := strings.Repeat("\x00", 256)

	gl.GetProgramiv(s.id, gl.ACTIVE_UNIFORMS, &count)
	for i := int32(0); i < count; i++ {
		gl.GetActiveUniform(s.id, uint32(i), int32(len(buf)), &length, &size, &tp, gl.Str(buf))

		// Force copy
		name := make([]byte, length)
		copy(name, []byte(buf[:length]))

		s.uniforms[string(name)] = gl.GetUniformLocation(s.id, gl.Str(string(name)+"\x00"))
	}
}

func loadShaderFromFiles(filenames ...string) (uint32, error) {
	data := make([]*ShaderData, 0, len(filenames))
	for _, file := range filenames {
		Loadf("asset.Shader [%v]", file)
		b, err := Load(file)
		if err != nil {
			return InvalidID, err
		}

		data = append(data, &ShaderData{
			Code: preProcessFile(file, string(b)),
			Type: getShaderType(file),
		})
	}
	return loadShaderFromData(data...)
}

func loadShaderFromData(data ...*ShaderData) (uint32, error) {
	pID := uint32(0)

	shaders := make([]uint32, 0, len(data))
	for _, d := range data {
		code := preProcessFile("", d.Code)
		id, err := compileShader(code, d.Type)
		if err != nil {
			Errorf("%v", err)
			Infof("Full Shader Code:\n%v", addLineNumbers(code))
			return 0, fmt.Errorf("Failed to compile shader")
		}
		shaders = append(shaders, id)
	}

	pID = gl.CreateProgram()
	for _, id := range shaders {
		gl.AttachShader(pID, id)
	}
	gl.LinkProgram(pID)

	var status int32
	gl.GetProgramiv(pID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(pID, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(pID, logLen, nil, gl.Str(log))

		return 0, fmt.Errorf("Failed to link program: %v", log)
	}

	for _, id := range shaders {
		gl.DeleteShader(id)
	}

	return pID, nil
}

func getVersionString() string {
	if _versionString != "" {
		return _versionString
	}

	tmp := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))

	_versionString = "#version "
	for i := range tmp {
		if tmp[i] == ' ' {
			break
		}
		if unicode.IsDigit(rune(tmp[i])) {
			_versionString += string(tmp[i])
		}
	}
	_versionString += " core\n"
	return _versionString
}

func preProcessFile(filename, code string) string {
	code = preProcessCode(filepath.Dir(filename), code, GetShaderDefines())

	// Prepend `#version`,
	code = getVersionString() + "\n" + code

	// Append null-terminator (windows)
	code += "\x00"

	return code
}

func preProcessCode(dir, code string, defines map[string]string) string {
	// Clean CRLF (windows)
	code = strings.Replace(code, "\r", "", -1)

	// PreProcessor statements
	lines := strings.Split(code+"\n", "\n")

	skipLines := false

	newLines := make([]string, 0, len(lines))
	for _, line := range lines {
		if len(line) > 0 && line[0] == '#' {
			if strings.HasPrefix(line, "#include") {
				file := strings.TrimSpace(line[8:])
				if len(file) < 3 {
					Warnf("No filename specified in #include")
					continue
				}
				first, last := file[0], file[len(file)-1]
				file = file[1 : len(file)-1]
				if first == '<' && last == '>' {
					file = filepath.Join(ShaderIncludePath, file)
				} else if first == '"' && last == '"' {
					file = filepath.Join(dir, file)
				} else {
					Warnf("Invalid #include format [%v]", line)
					Warnf("#include's must be either \"filename\" or <filename>")
					continue
				}

				b, err := Load(file)
				if err != nil {
					Warnf("Failed to include shader [%v]", file)
				}
				newLines = append(newLines,
					strings.Split(
						preProcessCode(file, string(b), defines),
						"\n")...)

			} else if strings.HasPrefix(line, "#define") {
				parts := strings.SplitN(strings.TrimSpace(line[8:]), " ", 2)
				name := parts[0]
				value := ""
				if len(parts) > 1 {
					value = parts[1]
				}
				defines[name] = value
			} else if strings.HasPrefix(line, "#undef") {
				name := strings.TrimSpace(line[8:])
				delete(defines, name)
			} else if strings.HasPrefix(line, "#ifdef") {
				name := strings.TrimSpace(line[7:])
				if _, found := defines[name]; !found {
					skipLines = true
				}
			} else if strings.HasPrefix(line, "#ifndef") {
				name := strings.TrimSpace(line[8:])
				if _, found := defines[name]; found {
					skipLines = true
				}
			} else if strings.HasPrefix(line, "#endif") {
				skipLines = false
			}
			continue
		}

		if skipLines {
			continue
		}

		for key, val := range defines {
			line = strings.Replace(line, key, val, -1)
		}

		newLines = append(newLines, line)
	}

	code = strings.Join(newLines, "\n")

	return code
}

func compileShader(code string, t uint32) (uint32, error) {
	id := gl.CreateShader(t)

	ccode, free := gl.Strs(code)
	gl.ShaderSource(id, 1, ccode, nil)
	free()
	gl.CompileShader(id)

	var status int32
	gl.GetShaderiv(id, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetShaderiv(id, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetShaderInfoLog(id, logLen, nil, gl.Str(log))

		return InvalidID, fmt.Errorf("Failed to compile shader: %v", log)
	}

	return id, nil
}

func getShaderType(filename string) uint32 {
	if strings.HasSuffix(filename, ".vs.glsl") {
		return gl.VERTEX_SHADER
	}
	if strings.HasSuffix(filename, ".fs.glsl") {
		return gl.FRAGMENT_SHADER
	}
	return gl.INVALID_ENUM
}
