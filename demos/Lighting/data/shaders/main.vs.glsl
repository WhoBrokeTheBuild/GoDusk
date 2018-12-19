#include <mvp.inc.glsl>
#include <attribute.inc.glsl>

uniform vec3 uPointLightPos;

out vec4 p_Position;
out vec4 p_Normal;
out vec2 p_TexCoord;

out vec4 p_LightPos;
// out vec4 p_LightDir;
// out vec4 p_ViewDir;

void main() {
    p_Position = uModel * vec4(_Position, 1.0);
    p_Normal   = vec4(mat3(transpose(inverse(uModel))) * _Normal, 1.0);
	p_TexCoord = vec2(_TexCoord.x, 1.0 - _TexCoord.y);
	
    p_LightPos = vec4(uPointLightPos, 1.0);
	// p_LightDir = normalize(-vec4(-0.2, -1.0, -0.3, 0.0));
	// p_ViewDir = -(uView * p_Position);

    gl_Position = uMVP * vec4(_Position, 1);
}
