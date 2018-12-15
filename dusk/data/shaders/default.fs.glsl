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
