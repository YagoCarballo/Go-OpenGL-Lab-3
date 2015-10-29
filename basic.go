package main
import (
	"fmt"
	"runtime"

	"./wrapper"
	"./objects"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const windowWidth = 1024
const windowHeight = 768
const windowFPS = 60

/* Define buffer object indices */
var positionBufferObject, colourObject, normalsBufferObject uint32

var shaderProgram uint32        /* Identifier for the shader prgoram */
var vertexArrayObject uint32            /* Vertex array (Containor) object. This is the index of the VAO that will be the container for
					   our buffer objects */

var colourmode objects.ColorMode    /* Index of a uniform to switch the colour mode in the vertex shader
					  I've included this to show you how to pass in an unsigned integer into
					  your vertex shader. */

/* Position and view globals */
var angle_x, angle_inc_x, x, scale, z, y float32
var angle_y, angle_inc_y, angle_z, angle_inc_z float32

var aspect_ratio float32        // Aspect ratio of the window defined in the reshape callback

// Uniforms
var modelUniform, projectionUniform, viewUniform int32
var colourmodeUniform int32

// Sphere
var sphere *objects.Sphere
var cube *objects.Cube


// Define vertices for a cube in 12 triangles
var vertexPositions = []float32{
	-0.25, 0.25, -0.25,
	-0.25, -0.25, -0.25,
	0.25, -0.25, -0.25,

	0.25, -0.25, -0.25,
	0.25, 0.25, -0.25,
	-0.25, 0.25, -0.25,

	0.25, -0.25, -0.25,
	0.25, -0.25, 0.25,
	0.25, 0.25, -0.25,

	0.25, -0.25, 0.25,
	0.25, 0.25, 0.25,
	0.25, 0.25, -0.25,

	0.25, -0.25, 0.25,
	-0.25, -0.25, 0.25,
	0.25, 0.25, 0.25,

	-0.25, -0.25, 0.25,
	-0.25, 0.25, 0.25,
	0.25, 0.25, 0.25,

	-0.25, -0.25, 0.25,
	-0.25, -0.25, -0.25,
	-0.25, 0.25, 0.25,

	-0.25, -0.25, -0.25,
	-0.25, 0.25, -0.25,
	-0.25, 0.25, 0.25,

	-0.25, -0.25, 0.25,
	0.25, -0.25, 0.25,
	0.25, -0.25, -0.25,

	0.25, -0.25, -0.25,
	-0.25, -0.25, -0.25,
	-0.25, -0.25, 0.25,

	-0.25, 0.25, -0.25,
	0.25, 0.25, -0.25,
	0.25, 0.25, 0.25,

	0.25, 0.25, 0.25,
	-0.25, 0.25, 0.25,
	-0.25, 0.25, -0.25,
}

// Define an array of colours
var vertexColours = []float32{
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

/* Manually specified normals for our cube */
var normals = []float32{
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	0, 0, -1,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	1, 0, 0,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	0, 0, 1,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	-1, 0, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, -1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
	0, 1, 0,
}

/////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////// Initialization ///////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

// This function is called by go as soon as this class is opened
func init() {
	// Locks the Execution in the main Thread as OpenGL is not thread Safe
	runtime.LockOSThread()
}

// Entry point of program
func main() {
	// Creates the Window Wrapper
	glw := wrapper.NewWrapper(windowWidth, windowHeight, "Lab 3: Lights")
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
	/* Set the object transformation controls to their initial values */
	x = 0.05;
	y = 0;
	z = 0;
	angle_x = 0
	angle_y = 0
	angle_z = 0
	angle_inc_x = 0
	angle_inc_y = 0
	angle_inc_z = 0
	scale = 1.0;
	aspect_ratio = 1.3333
	colourmode = objects.COLOR_SOLID
	var numLats uint32 = 20        // Number of latitudes in our sphere
	var numLongs uint32 = 20        // Number of longitudes in our sphere

	// Generate index (name) for one vertex array object
	gl.GenVertexArrays(1, &vertexArrayObject);

	// Create the vertex array object and make it current
	gl.BindVertexArray(vertexArrayObject);


	// Create the Cube Object
	cube = objects.NewCube(&vertexPositions, &vertexColours, &normals)
	cube.MakeVBO()

	// create the sphere object
	sphere = objects.NewSphere(numLats, numLongs);
	sphere.MakeSphereVBO()

	// Creates the Shader Program
	var err error; shaderProgram, err = wrapper.LoadShader("./shaders/basic.vert", "./shaders/basic.frag")

	// If there is any error loading the shaders, it panics
	if err != nil {
		panic(err)
	}

	// Define uniforms to send to vertex shader
	modelUniform = gl.GetUniformLocation(shaderProgram, gl.Str("model\x00"));
	colourmodeUniform = gl.GetUniformLocation(shaderProgram, gl.Str("colourmode\x00"));
	viewUniform = gl.GetUniformLocation(shaderProgram, gl.Str("view\x00"));
	projectionUniform = gl.GetUniformLocation(shaderProgram, gl.Str("projection\x00"));
}

/////////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////// Callbacks /////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////

//
// Draw Loop Function
// This function gets called on every update.
//
func drawLoop(glw *wrapper.Glw) {
	// Sets the Clear Color (Background Color)
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)

	// Clears the Window
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	// Enables Depth
	gl.Enable(gl.DEPTH_TEST)

	// Sets the Shader program to Use
	gl.UseProgram(shaderProgram)

	// Define the model transformations for the cube
	cube.ResetModel()
	cube.Translate(x + 0.5, y, z)
	cube.Scale(scale, scale, scale) //scale equally in all axis
	cube.Rotate(-angle_x, mgl32.Vec3{1, 0, 0}) //rotating in clockwise direction around x-axis
	cube.Rotate(-angle_y, mgl32.Vec3{0, 1, 0}) //rotating in clockwise direction around y-axis
	cube.Rotate(-angle_z, mgl32.Vec3{0, 0, 1}) //rotating in clockwise direction around z-axis

	// Define the model transformations for our sphere
	sphere.ResetModel()
	sphere.Translate(-x - 0.5, 0, 0)
	sphere.Scale(scale / 3.0, scale / 3.0, scale / 3.0) //scale equally in all axis
	sphere.Rotate(-angle_x, mgl32.Vec3{1, 0, 0}) //rotating in clockwise direction around x-axis
	sphere.Rotate(-angle_y, mgl32.Vec3{0, 1, 0}) //rotating in clockwise direction around y-axis
	sphere.Rotate(-angle_z, mgl32.Vec3{0, 0, 1}) //rotating in clockwise direction around z-axis

	// Projection matrix : 45° Field of View, 4:3 ratio, display range : 0.1 unit <-> 100 units
	var Projection mgl32.Mat4 = mgl32.Perspective(30.0, aspect_ratio, 0.1, 100.0)

	// Camera matrix
	var View mgl32.Mat4 = mgl32.LookAtV(
		mgl32.Vec3{0, 0, 4}, // Camera is at (0,0,4), in World Space
		mgl32.Vec3{0, 0, 0}, // and looks at the origin
		mgl32.Vec3{0, 1, 0}, // Head is up (set to 0,-1,0 to look upside-down)
	);

	// Send our uniforms variables to the currently bound shader,
	gl.Uniform1ui(colourmodeUniform, uint32(colourmode))
	gl.UniformMatrix4fv(viewUniform, 1, false, &View[0])
	gl.UniformMatrix4fv(projectionUniform, 1, false, &Projection[0])

	// Draws the Cube
	gl.UniformMatrix4fv(modelUniform, 1, false, &cube.Model[0])
	cube.Draw()

	// Draw our sphere
	gl.UniformMatrix4fv(modelUniform, 1, false, &sphere.Model[0])
	sphere.DrawSphere()

	gl.DisableVertexAttribArray(0);
	gl.UseProgram(0);

	/* Modify our animation variables */
	angle_x += angle_inc_x;
	angle_y += angle_inc_y;
	angle_z += angle_inc_z;
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
func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// React only if the key was just pressed
	if action != glfw.Press {
		return;
	}

	switch key {
	// If the Key Excape is pressed, it closes the App
	case glfw.KeyEscape:
		if action == glfw.Press {
			window.SetShouldClose(true)
		}
		break

	case glfw.KeyQ:
		angle_inc_x += 0.05
		break

	case glfw.KeyW:
		angle_inc_x -= 0.05
		break

	case glfw.KeyE:
		angle_inc_y += 0.05
		break

	case glfw.KeyR:
		angle_inc_y -= 0.05
		break

	case glfw.KeyT:
		angle_inc_z -= 0.05
		break

	case glfw.KeyY:
		angle_inc_z += 0.05
		break

	case glfw.KeyA:
		scale += 0.02
		break

	case glfw.KeyS:
		scale -= 0.02
		break

	case glfw.KeyZ:
		x -= 0.05
		break

	case glfw.KeyX:
		x += 0.05
		break

	case glfw.KeyC:
		y -= 0.05
		break

	case glfw.KeyV:
		y += 0.05
		break

	case glfw.KeyB:
		z -= 0.05
		break

	case glfw.KeyN:
		z += 0.05
		break

	case glfw.KeyM:
		if colourmode == 1 {
			colourmode = 0
		} else {
			colourmode = 1
		}
		fmt.Printf("Colour Mode: %s \n", colourmode)
		break

	// Cycle between drawing vertices, mesh and filled polygons
	case glfw.KeyK:
		sphere.DrawMode ++;
		if sphere.DrawMode > 2 {
			sphere.DrawMode = 0
		}
		break

	case glfw.KeyL:
		cube.DrawMode ++;
		if cube.DrawMode > 2 {
			cube.DrawMode = 0
		}
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
func reshape(window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height));
	aspect_ratio = (float32(width) / 640.0 * 4.0) / (float32(height) / 480.0 * 3.0);
}
