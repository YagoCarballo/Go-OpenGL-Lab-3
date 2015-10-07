package main
import (
	"runtime"
	"./wrapper"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"fmt"
//	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 1024
const windowHeight = 768
const windowFPS = 60

var shader uint32
var vertexArrayObject uint32

var vertexPositions = []float32{
	0.75, 0.75, 0.0, 1.0,
	0.75, -0.75, 0.0, 1.0,
	-0.75, -0.75, 0.0, 1.0,
}

func main() {
	// Locks the Execution in the main Thread
	runtime.LockOSThread()

	// Creates the Window Wrapper
	glw := wrapper.NewWrapper(windowWidth, windowHeight, "Hello Graphics World")
	glw.SetFPS(windowFPS)

	// Creates the Window
	glw.CreateWindow()

	// Creates the Vertex Array Object
	vertexArrayObject = createVertexArrayObject()

	// Sets the Viewport
	gl.Viewport(0, 0, int32(windowWidth), int32(windowHeight))

	// Sets the Event Callbacks
	glw.SetRenderCallback(drawLoop)
	glw.SetKeyCallBack(keyCallback)
	glw.SetReshapeCallback(reshape)

	shader = createShaderProgram()

	glw.StartLoop()
}

func move (key glfw.Key, action glfw.Action) {
	if key == glfw.KeyUp {
		vertexPositions[7] += 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeyDown {
		vertexPositions[7] -= 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeyLeft {
		vertexPositions[3] += 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeyRight {
		vertexPositions[3] -= 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeySpace {
		fmt.Println("Jump!!")
	}
}

// Callbacks

func drawLoop (glw *wrapper.Glw) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.BindVertexArray(vertexArrayObject)
	gl.UseProgram(shader)
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func keyCallback (window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	move(key, action)

	if action != glfw.Press {
		return;
	}

//	fmt.Println("Key Pressed: ", key)

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func reshape (window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

// Shaders

func createShaderProgram () uint32 {
	// Creates the Vertex Shader
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	cvertexShader := gl.Str(wrapper.FileToString("shaders/basic.vert"))
	gl.ShaderSource(vertexShader, 1, &cvertexShader, nil)
	gl.CompileShader(vertexShader)

	// Creates the Fragment Shader
	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	cfragmentShader := gl.Str(wrapper.FileToString("shaders/basic.frag"))
	gl.ShaderSource(fragmentShader, 1, &cfragmentShader, nil)
	gl.CompileShader(fragmentShader)

	// Creates the Shader Program and attaches the shaders to it
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.AttachShader(shaderProgram, vertexShader)

	// Links the Shader Program
	gl.LinkProgram(shaderProgram)

	return shaderProgram
}

// 3D

func createVertexBufferObject (points []float32) uint32 {
	var vertexBufferObject uint32
	gl.GenBuffers(1, &vertexBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(points) * 4, gl.Ptr(&points[0]), gl.STATIC_DRAW)
	return vertexBufferObject
}

func createVertexArrayObject () uint32 {
	var vertexBufferObject = createVertexBufferObject(vertexPositions)

	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	return vertexArrayObject
}
