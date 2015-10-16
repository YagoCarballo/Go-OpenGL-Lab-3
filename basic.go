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

// Define vertex positions in homogeneous coordinates
var vertexPositions = []float32{
	0.75, 0.75, 0.0, 1.0,
	0.75, -0.75, 0.0, 1.0,
	-0.75, -0.75, 0.0, 1.0,
}

// Define an array of colors
var vertexColours = []float32 {
	1.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
};

func init () {
	// Locks the Execution in the main Thread
	runtime.LockOSThread()
}

func main() {
	// Creates the Window Wrapper
	glw := wrapper.NewWrapper(windowWidth, windowHeight, "Hello Graphics World")
	glw.SetFPS(windowFPS)

	// Creates the Window
	glw.CreateWindow()

	// Creates the Vertex Array Object
	vertexArrayObject = createVertexArrayObject()

	// Sets the Event Callbacks
	glw.SetRenderCallback(drawLoop)
	glw.SetKeyCallBack(keyCallback)
	glw.SetReshapeCallback(reshape)

	// Sets the Viewport (Important !!, this has to run at the end!!)
	defer gl.Viewport(0, 0, windowWidth, windowHeight)

	// Creates the Shader Program
	shader = createShaderProgram()

	// Starts the Rendering Loop
	glw.StartLoop()
}

func move (key glfw.Key, action glfw.Action) {
	if key == glfw.KeyUp {
		vertexPositions[1] += 0.05
		vertexPositions[5] += 0.05
		vertexPositions[9] += 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeyDown {
		vertexPositions[1] -= 0.05
		vertexPositions[5] -= 0.05
		vertexPositions[9] -= 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeyLeft {
		vertexPositions[0] -= 0.05
		vertexPositions[4] -= 0.05
		vertexPositions[8] -= 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeyRight {
		vertexPositions[0] += 0.05
		vertexPositions[4] += 0.05
		vertexPositions[8] += 0.05
		vertexArrayObject = createVertexArrayObject()

	} else if key == glfw.KeySpace {
		fmt.Println("Jump!!")
	}
}

// Callbacks

// Draw Loop Function
// This function gets called on every update.
func drawLoop (glw *wrapper.Glw) {
	// Sets the Clear Color (Background Color)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// Clears the Window
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Loads the Shader
	gl.UseProgram(shader)

	// Enables the generic vertex attribute array
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	// Binds the vertex array object
	gl.BindVertexArray(vertexArrayObject)

	// Draws the Vertex Array
	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	// Disables the generic vertex attribute array
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)

	// Disables the Shaders
	gl.UseProgram(0)
}

// This function gets called when a key is pressed
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

// Creates the Vertex Array Object
// Creates the buffer, then creates the vertex array, and binds the buffer to it.
func createVertexArrayObject () uint32 {
	// Creates the Vertex Buffer
	var vertexBufferObject = createVertexBufferObject(vertexPositions)

	// Genrate a name for a vertex buffer object
	var colourObject = createVertexBufferObject(vertexColours)

	// Creates the Buffer Array
	var vertexArrayObject uint32
	gl.GenVertexArrays(1, &vertexArrayObject)
	gl.BindVertexArray(vertexArrayObject)

	// Binds Positions Buffer to the Vertex Array?
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexBufferObject)
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil)
	gl.DisableVertexAttribArray(1)

	// Binds Color Buffer to the Vertex Array?
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, colourObject)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)
	gl.DisableVertexAttribArray(0)


	return vertexArrayObject
}
