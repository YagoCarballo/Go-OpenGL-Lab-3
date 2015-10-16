package main
import (
	"runtime"

	"./wrapper"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 1024
const windowHeight = 768
const windowFPS = 60

var positionBufferObject, colourObject uint32
var shaderProgram uint32
var vertexArrayObject uint32

// Position and view globals
var angle_x, angle_x_inc float64
var angle_y, angle_y_inc float64
var angle_z, angle_z_inc float64
var scale float32 = 1.0

// Uniforms
var modelID int32

var model mgl32.Mat4

// Define vertices for a cube in 12 triangles
var vertexPositions = []float32{
	-0.25, 0.25, -0.25, 1.0,
	-0.25, -0.25, -0.25, 1.0,
	0.25, -0.25, -0.25, 1.0,

	0.25, -0.25, -0.25, 1.0,
	0.25, 0.25, -0.25, 1.0,
	-0.25, 0.25, -0.25, 1.0,

	0.25, -0.25, -0.25, 1.0,
	0.25, -0.25, 0.25, 1.0,
	0.25, 0.25, -0.25, 1.0,

	0.25, -0.25, 0.25, 1.0,
	0.25, 0.25, 0.25, 1.0,
	0.25, 0.25, -0.25, 1.0,

	0.25, -0.25, 0.25, 1.0,
	-0.25, -0.25, 0.25, 1.0,
	0.25, 0.25, 0.25, 1.0,

	-0.25, -0.25, 0.25, 1.0,
	-0.25, 0.25, 0.25, 1.0,
	0.25, 0.25, 0.25, 1.0,

	-0.25, -0.25, 0.25, 1.0,
	-0.25, -0.25, -0.25, 1.0,
	-0.25, 0.25, 0.25, 1.0,

	-0.25, -0.25, -0.25, 1.0,
	-0.25, 0.25, -0.25, 1.0,
	-0.25, 0.25, 0.25, 1.0,

	-0.25, -0.25, 0.25, 1.0,
	0.25, -0.25, 0.25, 1.0,
	0.25, -0.25, -0.25, 1.0,

	0.25, -0.25, -0.25, 1.0,
	-0.25, -0.25, -0.25, 1.0,
	-0.25, -0.25, 0.25, 1.0,

	-0.25, 0.25, -0.25, 1.0,
	0.25, 0.25, -0.25, 1.0,
	0.25, 0.25, 0.25, 1.0,

	0.25, 0.25, 0.25, 1.0,
	-0.25, 0.25, 0.25, 1.0,
	-0.25, 0.25, -0.25, 1.0,
}

// Define an array of colours
var vertexColours = []float32 {
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,
	0.0, 0.0, 1.0, 1.0,

	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,
	0.0, 1.0, 0.0, 1.0,

	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 0.0, 1.0,

	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,
	1.0, 0.0, 0.0, 1.0,

	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,
	1.0, 0.0, 1.0, 1.0,

	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
	0.0, 1.0, 1.0, 1.0,
}

/////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// Initialization ///////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

// This function is called by go as soon as this class is opened
func init () {
	// Locks the Execution in the main Thread as OpenGL is not thread Safe
	runtime.LockOSThread()
}

// Entry point of program
func main() {
	// Creates the Window Wrapper
	glw := wrapper.NewWrapper(windowWidth, windowHeight, "Hello Graphics World")
	glw.SetFPS(windowFPS)

	// Creates the Window
	glw.CreateWindow()

	// Sets the Event Callbacks
	glw.SetRenderCallback(drawLoop)
	glw.SetKeyCallBack(keyCallback)
	glw.SetReshapeCallback(reshape)

	// Initializes the App
	InitApp(glw)

	// Starts the Rendering Loop
	glw.StartLoop()

	// Sets the Viewport (Important !!, this has to run after the loop!!)
	defer gl.Viewport(0, 0, windowWidth, windowHeight)
}

//
// Init App
// This function initializes the variables and sets up the environment.
//
// @param wrapper (*wrapper.Glw) the window wrapper
//
func InitApp(glw *wrapper.Glw) {
	angle_x = 0;
	angle_x_inc = 0;

	// Initializes a model in the start position
	model = mgl32.Ident4()

	// Generate index (name) for one vertex array object
	gl.GenVertexArrays(1, &vertexArrayObject);

	// Create the vertex array object and make it current
	gl.BindVertexArray(vertexArrayObject);

	// Create a vertex buffer object to store vertices
	gl.GenBuffers(1, &positionBufferObject);
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject);
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexPositions) * 4, gl.Ptr(vertexPositions), gl.STATIC_DRAW);
	gl.BindBuffer(gl.ARRAY_BUFFER, 0);

	// Create a vertex buffer object to store vertex colours
	gl.GenBuffers(1, &colourObject);
	gl.BindBuffer(gl.ARRAY_BUFFER, colourObject);
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexColours) * 4, gl.Ptr(vertexColours), gl.STATIC_DRAW);
	gl.BindBuffer(gl.ARRAY_BUFFER, 0);

	// Creates the Shader Program
	var err error; shaderProgram, err = wrapper.LoadShader("./shaders/basic.vert", "./shaders/basic.frag")

	// If there is any error loading the shaders, it panics
	if err != nil {
		panic(err)
	}

	// Define uniforms to send to vertex shader
	modelID = gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"));
}

/////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////// Callbacks /////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

//
// Draw Loop Function
// This function gets called on every update.
//
func drawLoop (glw *wrapper.Glw) {
	// Sets the Clear Color (Background Color)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// Clears the Window
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Sets the Shader program to Use
	gl.UseProgram(shaderProgram)

	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject);
	gl.EnableVertexAttribArray(0);

	// glVertexAttribPointer(index, size, type, normalised, stride, pointer)
	// index relates to the layout qualifier in the vertex shader and in
	// glEnableVertexAttribArray() and glDisableVertexAttribArray()
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil);

	gl.BindBuffer(gl.ARRAY_BUFFER, colourObject);
	gl.EnableVertexAttribArray(1);

	// glVertexAttribPointer(index, size, type, normalised, stride, pointer)
	// index relates to the layout qualifier in the vertex shader and in
	// glEnableVertexAttribArray() and glDisableVertexAttribArray()
	gl.VertexAttribPointer(1, 4, gl.FLOAT, false, 0, nil);

	// Rotates the model
	modelX := mgl32.HomogRotate3D(float32(angle_y), mgl32.Vec3{1, 0, 0})
	modelY := mgl32.HomogRotate3D(float32(angle_x), mgl32.Vec3{0, 1, 0})
	modelZ := mgl32.HomogRotate3D(float32(angle_z), mgl32.Vec3{0, 0, 1})
	modelScale := mgl32.Scale3D(scale, scale, scale)

	// Multiplies both cubes to apply both rotations
	model = modelX.Mul4(modelY).Mul4(modelZ).Mul4(modelScale)

	// Send our transformations to the currently bound shader
	gl.UniformMatrix4fv(modelID, 1, false, &model[0])

	// Draws the Cube
	gl.DrawArrays(gl.TRIANGLES, 0, 36)

	gl.DisableVertexAttribArray(0);

	// Disables the Shaders
	gl.UseProgram(0)

	// Modify our animation variables
	angle_x += angle_x_inc;
	angle_y += angle_y_inc;
	angle_z += angle_z_inc;
}

//
// key Callback
// This function gets called when a key is pressed
//
// @param window (*glfw.Window) a pointer to the window
// @param key (glfw.Key) the pressed key
// @param scancode (int) the scancode
// @param action (glfw.Action) the state of the key
// @param mods (glfw.ModifierKey) the pressed modified keys.
//
func keyCallback (window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// React only if the key was just pressed
	if action != glfw.Press {
		return;
	}

	// If the Key Excape is pressed, it closes the App
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}

	// If the Key W is pressed, it rotates up
	if key == glfw.KeyW || key == glfw.KeyUp {
		angle_y_inc += 0.1
	}

	// If the Key Q is pressed, it rotates down
	if key == glfw.KeyS || key == glfw.KeyDown {
		angle_y_inc -= 0.1
	}

	// If the Key A is pressed, it rotates to the Left
	if key == glfw.KeyA || key == glfw.KeyLeft {
		angle_x_inc -= 0.1
	}

	// If the Key D is pressed, it rotates to the Right
	if key == glfw.KeyD || key == glfw.KeyRight {
		angle_x_inc += 0.1
	}

	// If the Key Q is pressed, it rotates to the Back
	if key == glfw.KeyQ {
		angle_z_inc -= 0.1
	}

	// If the Key E is pressed, it rotates to the Front
	if key == glfw.KeyE {
		angle_z_inc += 0.1
	}

	// If the Key Z is pressed, it Scales Out
	if key == glfw.KeyZ {
		scale -= 0.5
	}

	// If the Key X is pressed, it Scales In
	if key == glfw.KeyX {
		scale += 0.5
	}
}

//
// Reshape
// This gets called when the window changes its size
//
// @param window (*glfw.Window) a pointer to the window
// @param width (int) the width of the window
// @param height (int) the height of the window
//
func reshape (window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}
