package dusk

import gl "github.com/go-gl/gl/v4.1-core/gl"

const (
	uiShaderVert = `
#include <mvp.inc.glsl>
#include <attribute.inc.glsl>

out vec2 p_TexCoord;

void main() {
    p_TexCoord = vec2(_TexCoord.x, 1.0 - _TexCoord.y);

    gl_Position = uProjection * vec4(_Position, 1);
}
`
	uiShaderFrag = `
uniform sampler2D uTexture; 

in vec2 p_TexCoord;

out vec4 _Color;

void main() {
	_Color = texture(uTexture, p_TexCoord);
}
`
)

// UIShader represents the default UI shader
type UIShader struct {
	Shader
}

var _uiShader *UIShader

// GetUIShader returns an instance of the UIShader
func GetUIShader() *UIShader {
	if _uiShader != nil {
		return _uiShader
	}
	_uiShader = &UIShader{}
	_uiShader.InitFromData(
		&ShaderData{
			Code: uiShaderVert,
			Type: gl.VERTEX_SHADER,
		},
		&ShaderData{
			Code: uiShaderFrag,
			Type: gl.FRAGMENT_SHADER,
		},
	)
	return _uiShader
}

// Bind implements the Shader interface
func (s *UIShader) Bind(ctx *RenderContext, data interface{}) {
	s.Shader.Bind(ctx, data)

	gl.UniformMatrix4fv(s.UniformLocation("uProjection"), 1, false, &ctx.Projection[0])
}
