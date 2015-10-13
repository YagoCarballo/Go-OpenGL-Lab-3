// Minimal vertex shader

#version 410
layout(location = 0) in vec4 pos;
out vec4 position;

void main()
{
    position = pos;
	gl_Position = pos;
}
