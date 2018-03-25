#version 330 core

layout (location = 0) in vec3 pos;

uniform vec2 center;
uniform float scale;

out vec2 position;

void main() {
  gl_Position = vec4(pos.xy, 0.0f, 1.0f);
  position = (center + (pos.xy * scale));
}
