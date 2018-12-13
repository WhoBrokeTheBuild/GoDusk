uniform vec4 uAmbient;
uniform vec4 uDiffuse;
uniform vec4 uSpecular;

uniform sampler2D uAmbientMap; 
uniform sampler2D uDiffuseMap; 
uniform sampler2D uSpecularMap; 

in vec4 p_Position;
in vec4 p_Normal;
in vec2 p_TexCoord;

in vec3 p_LightDir;
in vec3 p_ViewDir;

out vec4 _Color;

void main() {
    vec3 normal = normalize(p_Normal.xyz);
    
    float d = max(dot(normal, p_LightDir), 0.0);
    vec4 diffuse = vec4(d, d, d, 1.0) * (uDiffuse + texture(uDiffuseMap, p_TexCoord));

    _Color = diffuse;
}
