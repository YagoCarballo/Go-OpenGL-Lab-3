// Minimal vertex shader

#version 330
layout(location = 0) in vec4 position;
layout(location = 1) in vec4 colour;
out vec4 fcolour;
uniform mat4 model;
uniform mat4 projection;
uniform mat4 camera;

void main()
{
	gl_Position = projection * camera * model * position;

	fcolour = colour;
//	fcolour = position * 2.0 + vec4(0.5, 0.5, 0.5, 1.0);
}