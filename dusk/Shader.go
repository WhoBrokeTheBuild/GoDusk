package dusk

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	gl "github.com/go-gl/gl/v4.1-core/gl"
)

const (
	// ShaderIncludePath is the path to search for files from #include <>
	ShaderIncludePath = "data/shaders/include"
)

// Shader represents an OpenGL Shader Program
type Shader struct {
	ID       uint32
	Uniforms map[string]int32
}

var _versionString string
var _defaultDefines = map[string]string{}

func RegisterShaderDefines(defines map[string]interface{}) {
	for k, v := range defines {
		_defaultDefines[k] = fmt.Sprintf("%v", v)
	}
}

func GetShaderDefines() map[string]string {
	tmp := map[string]string{}
	for k, v := range _defaultDefines {
		tmp[k] = v
	}
	return tmp
}

// NewShaderFromFiles returns a new Shader from the given files
func NewShaderFromFiles(filenames []string) (*Shader, error) {
	s := &Shader{
		ID:       InvalidID,
		Uniforms: map[string]int32{},
	}

	err := s.LoadFromFiles(filenames)
	if err != nil {
		s.Delete()
		return nil, err
	}

	return s, nil
}

// Delete frees all resources owned by the Shader
func (s *Shader) Delete() {
	if s.ID != InvalidID {
		gl.DeleteProgram(s.ID)
		s.ID = InvalidID
	}
}

// LoadFromFiles loads a shader from the given files
func (s *Shader) LoadFromFiles(filenames []string) error {
	s.Delete()

	shaders := make([]uint32, 0, len(filenames))
	for _, file := range filenames {
		id, err := compileShader(file)
		if err != nil {
			return err
		}
		shaders = append(shaders, id)
	}

	s.ID = gl.CreateProgram()
	for _, id := range shaders {
		gl.AttachShader(s.ID, id)
	}
	gl.LinkProgram(s.ID)

	var status int32
	gl.GetProgramiv(s.ID, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLen int32
		gl.GetProgramiv(s.ID, gl.INFO_LOG_LENGTH, &logLen)

		log := strings.Repeat("\x00", int(logLen+1))
		gl.GetProgramInfoLog(s.ID, logLen, nil, gl.Str(log))

		s.Delete()
		return fmt.Errorf("Failed to link program: %v", log)
	}

	for _, id := range shaders {
		gl.DeleteShader(id)
	}

	s.cacheUniforms()

	return nil
}

// Bind calls glUseProgram with this Shader's ID
func (s *Shader) Bind() error {
	if s.ID == InvalidID {
		return fmt.Errorf("Failed to bind program: Not loaded")
	}

	gl.UseProgram(s.ID)
	return nil
}

// GetUniformLocation returns the uniform's location ID, or -1
func (s *Shader) GetUniformLocation(name string) int32 {
	if u, ok := s.Uniforms[name]; ok {
		return u
	}
	return -1
}

func (s *Shader) cacheUniforms() {
	var count int32

	var size int32
	var length int32
	var tp uint32
	buf := strings.Repeat("\x00", 256)

	gl.GetProgramiv(s.ID, gl.ACTIVE_UNIFORMS, &count)
	for i := int32(0); i < count; i++ {
		gl.GetActiveUniform(s.ID, uint32(i), int32(len(buf)), &length, &size, &tp, gl.Str(buf))

		// Force copy
		name := make([]byte, length)
		copy(name, []byte(buf[:length]))

		s.Uniforms[string(name)] = gl.GetUniformLocation(s.ID, gl.Str(string(name)+"\x00"))
	}
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
	code = preProcessCode(filename, code, GetShaderDefines())

	// Prepend `#version`,
	code = getVersionString() + "\n" + code

	// Append null-terminator (windows)
	code += "\x00"

	return code
}

func preProcessCode(filename, code string, defines map[string]string) string {
	dir := filepath.Dir(filename)

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

func compileShader(filename string) (uint32, error) {
	filename = filepath.Clean(filename)

	t := getShaderType(filename)
	id := gl.CreateShader(t)

	Loadf("asset.Shader [%v]", filename)
	b, err := Load(filename)
	if err != nil {
		return InvalidID, err
	}

	code := preProcessFile(filename, string(b))

	re := regexp.MustCompile(`\r`)
	code = re.ReplaceAllString(code, "")

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

		Infof("Full Shader Code:\n%v", addLineNumbers(code))

		return InvalidID, fmt.Errorf("Failed to compile [%v]: %v", filename, log)
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
