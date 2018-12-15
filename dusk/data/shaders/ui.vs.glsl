#include <mvp.inc.glsl>
#include <attribute.inc.glsl>

out vec2 p_TexCoord;

void main() {
    p_TexCoord = vec2(_TexCoord.x, 1.0 - _TexCoord.y);

    gl_Position = uProjection * vec4(_Position, 1);
}
