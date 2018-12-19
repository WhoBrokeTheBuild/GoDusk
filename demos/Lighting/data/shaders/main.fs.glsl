#include <material.inc.glsl>

in vec4 p_Position;
in vec4 p_Normal;
in vec2 p_TexCoord;

in vec4 p_LightPos;
// in vec4 p_LightDir;
// in vec4 p_ViewDir;

out vec4 _Color;

void main() {
    vec4 normal = p_Normal;
    if (HasNormalMap()) {
        normal = texture(uNormalMap, p_TexCoord);
    }
    normal = normalize(normal);

    vec4 ambient = uAmbient;
    if (HasAmbientMap()) {
        ambient = texture(uAmbientMap, p_TexCoord);
    }
    ambient *= 0.1;

    vec4 lightDir = normalize(p_LightPos - p_Position);
    float diff = max(0.0, dot(normal.xyz, lightDir.xyz));

    vec4 diffuse = uDiffuse;
    if (HasDiffuseMap()) {
        diffuse = texture(uDiffuseMap, p_TexCoord);
    }
    diffuse = vec4(diff * diffuse.rgb, diffuse.a);

    float spec = max(0.0, dot(normal.xyz, lightDir.xyz));
    spec *= spec;

    vec4 specular = uSpecular;
    if (HasSpecularMap()) {
        specular = texture(uSpecularMap, p_TexCoord);
    }
    specular = vec4(spec * vec3(specular.rgb), specular.a);

    _Color = vec4(ambient.rgb + diffuse.rgb + specular.rgb, diffuse.a);
}
