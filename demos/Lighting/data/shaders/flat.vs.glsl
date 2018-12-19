#include <mvp.inc.glsl>
#include <attribute.inc.glsl>

void main() {
    gl_Position = uMVP * vec4(_Position, 1);
}
