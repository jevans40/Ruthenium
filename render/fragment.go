package render

var fragmentShader = `

#version 450 core

uniform sampler2D texture1;
//uniform sampler2D texture2;

in vec2 TexPos;
in vec4 Color;
in flat int TexMap;

out vec4 FragColor;

void main(){
  if (TexMap > 0){
    FragColor = texture(texture1,TexPos) * Color;
  }
  else {
    FragColor = Color * vec4(1,1,1,1);
  }
}
` + "\x00"
