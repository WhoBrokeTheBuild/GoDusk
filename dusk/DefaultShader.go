package dusk

import (
	gl "github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	defaultShaderVert = `
#include <mvp.inc.glsl>
#include <attribute.inc.glsl>

uniform vec3 uCamera;

out vec4 p_Position;
out vec4 p_Normal;
out vec2 p_TexCoord;

out vec3 p_LightDir;
out vec3 p_ViewDir;

void main() {
    p_Position = uModel * vec4(_Position, 1.0);
    p_Normal   = uModel * vec4(_Normal, 1.0);
    p_TexCoord = vec2(_TexCoord.x, 1.0 - _TexCoord.y);

    p_LightDir = normalize(vec3(1000, 1000, 1000) - p_Position.xyz);
    p_ViewDir  = normalize(uCamera - p_Position.xyz);

    gl_Position = uProjection * uView * uModel * vec4(_Position, 1);
}
`
	defaultShaderFrag = `
#include <material.inc.glsl>

in vec4 p_Position;
in vec4 p_Normal;
in vec2 p_TexCoord;

in vec3 p_LightDir;
in vec3 p_ViewDir;

out vec4 _Color;

void main() {
    vec4 normal = normalize(p_Normal);
    if (HasNormalMap()) {
        normal = texture(uNormalMap, p_TexCoord);
    }

    vec4 ambient = uAmbient;
    if (HasAmbientMap()) {
        ambient = texture(uAmbientMap, p_TexCoord);
    }
    ambient *= 0.1;

    float diff = max(0.0, dot(normal.xyz, p_LightDir));

    vec4 diffuse = uDiffuse;
    if (HasDiffuseMap()) {
        diffuse = texture(uDiffuseMap, p_TexCoord);
    }
    diffuse = vec4(diff * vec3(diffuse.rgb), diffuse.a);

    vec3  half = normalize(p_LightDir + p_ViewDir);
    float spec = pow(max(0.0, dot(normal.xyz, half)), 32.0);

    vec4 specular = uSpecular;
    if (HasSpecularMap()) {
        specular = texture(uSpecularMap, p_TexCoord);
    }
    specular = vec4(spec * vec3(specular.rgb), specular.a);

    _Color = vec4(ambient.rgb + diffuse.rgb + specular.rgb, diffuse.a);
}
	`
)

// DefaultShader is the default shader used to render meshes
type DefaultShader struct {
	Shader
}

var _defaultShader *DefaultShader

// GetDefaultShader returns an instance of the DefaultShader
func GetDefaultShader() *DefaultShader {
	if _defaultShader != nil {
		return _defaultShader
	}
	_defaultShader = &DefaultShader{}
	_defaultShader.InitFromData(
		&ShaderData{
			Code: defaultShaderVert,
			Type: gl.VERTEX_SHADER,
		},
		&ShaderData{
			Code: defaultShaderFrag,
			Type: gl.FRAGMENT_SHADER,
		},
	)
	return _defaultShader
}

// Bind implements the Shader interface
func (s *DefaultShader) Bind(ctx *RenderContext, data interface{}) {
	s.Shader.Bind(ctx, data)
	model := mgl32.Mat4{}
	if data != nil {
		model = data.(mgl32.Mat4)
	}

	gl.UniformMatrix4fv(s.UniformLocation("uProjection"), 1, false, &ctx.Projection[0])
	gl.UniformMatrix4fv(s.UniformLocation("uView"), 1, false, &ctx.Camera.View[0])
	gl.UniformMatrix4fv(s.UniformLocation("uModel"), 1, false, &model[0])
	gl.Uniform4fv(s.UniformLocation("uCamera"), 1, &ctx.Camera.Position[0])
}
