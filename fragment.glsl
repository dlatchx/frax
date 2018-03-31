#version 330 core

#define cplx_add(a, b) vec2(a.x+b.x, a.y+b.y)
#define cplx_mul(a, b) vec2(a.x*b.x-a.y*b.y, a.x*b.y+a.y*b.x)
#define cplx_div(a, b) vec2(((a.x*b.x+a.y*b.y)/(b.x*b.x+b.y*b.y)), ((a.y*b.x-a.x*b.y)/(b.x*b.x+b.y*b.y)))
#define cplx_conj(a)   vec2(a.x, -a.y)
#define cplx_exp(z)    vec2(exp(z.x)*cos(z.y), exp(z.x)*sin(z.y))

uniform float scale;
uniform vec2 c;
uniform int max_it;

in vec2 position;

vec3 hsv2rgb(vec3 c) {
  vec4 K = vec4(1.0, 2.0 / 3.0, 1.0 / 3.0, 3.0);
  vec3 p = abs(fract(c.xxx + K.xyz) * 6.0 - K.www);
  return c.z * mix(K.xxx, clamp(p - K.xxx, 0.0, 1.0), c.y);
}

void main() {
  vec2 z = position;
  int i = 0;
  while (i < max_it && (z.x * z.x + z.y * z.y) < 4.0f) {
    z = cplx_add(cplx_mul(z, z), c);
    i++;
  }
  float it = float(i) / float(max_it);

  it += 246.0f / 360.0f;
  if (it > 1.0f)
    it -= 1.0f;
  vec3 hsv = vec3(it, 1.0f, 1.0f);
  gl_FragColor.xyz = hsv2rgb(hsv);
  gl_FragColor.w = 1.0f;
}
