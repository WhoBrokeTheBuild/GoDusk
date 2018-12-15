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
