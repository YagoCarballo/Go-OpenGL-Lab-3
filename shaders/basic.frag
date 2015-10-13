// Minimal fragment shader

#version 410

in vec4 position;
out vec4 outputColor;
void main()
{
	outputColor = vec4(vec3(position.x, position.y, position.z) * 0.5f, 1.0f);
//	outputColor = vec4(vec3(position.x, position.y, position.z) * 0.5f, 1.0f);
//	outputColor = vec4(1.0f, 1.0f, 0.0f, 1.0f);
}
