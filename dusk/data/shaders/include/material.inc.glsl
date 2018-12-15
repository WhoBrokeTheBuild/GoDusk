#ifndef MATERIAL_INC
#define MATERIAL_INC

uniform vec4 uAmbient;
uniform vec4 uDiffuse;
uniform vec4 uSpecular;

uniform sampler2D uAmbientMap; 
uniform sampler2D uDiffuseMap; 
uniform sampler2D uSpecularMap; 
uniform sampler2D uNormalMap; 

uniform uint uMapFlags;

bool HasAmbientMap() {
    return ((uMapFlags & FLAG_AMBIENT_MAPu) > 0u);
}
bool HasDiffuseMap() {
    return ((uMapFlags & FLAG_DIFFUSE_MAPu) > 0u);
}
bool HasSpecularMap() {
    return ((uMapFlags & FLAG_SPECULAR_MAPu) > 0u);
}
bool HasNormalMap() {
    return ((uMapFlags & FLAG_NORMAL_MAPu) > 0u);
}

#endif MATERIAL_INC