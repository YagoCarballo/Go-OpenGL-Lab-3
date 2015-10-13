// Minimal vertex shader

#version 410
layout(location = 0) in vec4 position;
layout(location = 1) in vec4 colour;
out vec4 fcolour;

void main()
{
	gl_Position = position;
	fcolour = colour;
}