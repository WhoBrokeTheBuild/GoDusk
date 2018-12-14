uniform vec4 uAmbient;
uniform vec4 uDiffuse;
uniform vec4 uSpecular;

uniform sampler2D uAmbientMap; 
uniform sampler2D uDiffuseMap; 
uniform sampler2D uSpecularMap; 
uniform sampler2D uBumpMap; 

in vec4 p_Position;
in vec4 p_Normal;
in vec2 p_TexCoord;

in vec3 p_LightDir;
in vec3 p_ViewDir;

out vec4 _Color;

void main() {
    vec4 normal = texture(uBumpMap, p_TexCoord);
    if (normal.w == 0) {
        normal = normalize(p_Normal);
    }

    vec4 ambient = uAmbient + texture(uAmbientMap, p_TexCoord);

    vec4 diffuse = uDiffuse + texture(uDiffuseMap, p_TexCoord);
    float diff = max(0.0, dot(normal.xyz, p_LightDir));
    diffuse = vec4(diff * vec3(diffuse.rgb), diffuse.a);

    vec3  half = normalize(p_LightDir + p_ViewDir);
    float spec = pow(max(0.0, dot(normal.xyz, half)), 32.0);

    vec4 specular = spec * (uSpecular + texture(uSpecularMap, p_TexCoord));

    _Color = vec4(ambient.rgb + diffuse.rgb + specular.rgb, diffuse.a);
}
