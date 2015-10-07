package wrapper

import (
	"log"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"fmt"
	"strings"
	"io/ioutil"
	"runtime"
)

type Glw struct  {
	Width, Height int
	Title string
	fps int
	running bool
	Window *glfw.Window
	renderer func()
	keyCallBack func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey)
	reshape func(*glfw.Window, int, int)
}

func NewGlw(width, height int, title string) *Glw {
	runtime.LockOSThread()

	// Init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	return &Glw{ width, height, title, 60, true, nil, nil, nil, nil }
}

func (glw *Glw) OpenWindow () error {
	glfw.WindowHint(glfw.Samples, 4)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)    // Necessary for OS X
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // Necessary for OS X

	// Debug
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	var err error = nil
	// Create a GLFW window, bail out if it doesn't work
	glw.Window, err = glfw.CreateWindow(glw.Width, glw.Height, glw.Title, nil, nil)
	if err != nil {
		glfw.Terminate();
		return err
	}

	// Obtain an OpenGL context and assign to the just opened GLFW window
	glw.Window.MakeContextCurrent();
	glw.Window.SetTitle(glw.Title)

	glw.Window.SetInputMode(glfw.StickyKeysMode, 1)

	return nil
}

func (glw *Glw) LoadShader (vertex_path string, fragment_path string) (uint32, error) {
	var vertShader, fragShader uint32;

	// Read shaders
	vertShaderStr, err := readFile(vertex_path)
	if err != nil {
		return 0, err
	}

	fragShaderStr, err := readFile(fragment_path)
	if err != nil {
		return 0, err
	}

	var result int32 = gl.FALSE
	var logLength int32

	vertShader = BuildShader(gl.VERTEX_SHADER, vertShaderStr);
	fragShader = BuildShader(gl.FRAGMENT_SHADER, fragShaderStr);

	fmt.Println("Linking program")
	program := gl.CreateProgram();
	gl.AttachShader(program, vertShader);
	gl.AttachShader(program, fragShader);
	gl.LinkProgram(program);

	gl.GetProgramiv(program, gl.LINK_STATUS, &result);
	gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength);

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))
	fmt.Println(log)

	gl.DeleteShader(vertShader);
	gl.DeleteShader(fragShader);

	return program, nil
}

func (glw *Glw) SetRenderer(display func()) {
	glw.renderer = display
}

func (glw *Glw) SetKeyCallback(keyCallback func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey)) {
	glw.keyCallBack = keyCallback
	glw.Window.SetKeyCallback(keyCallback)
}

func (glw *Glw) SetReshapeCallback(reshape func(*glfw.Window, int, int)) {
	glw.reshape = reshape
	glw.Window.SetFramebufferSizeCallback(reshape)
}

func (glw *Glw) EventLoop() {
	/* The event loop */
	for !glw.Window.ShouldClose() {
		// Call function to draw your graphics
		glw.renderer()

		// Swap buffers
		glw.Window.SwapBuffers()
		glfw.PollEvents()
	}

	glw.Terminate();
}

func (glw *Glw) Terminate () {
	/* Clean up */
	glw.Window.Destroy()
	glfw.Terminate()
}

func readFile (path string) (string, error) {
	fileContents, err := ioutil.ReadFile(path)
	return string(fileContents) + "\x00", err
}

func BuildShader(eShaderType uint32, shaderText string) uint32 {
	shader := gl.CreateShader(eShaderType)
	strFileData := gl.Str(shaderText)
	gl.ShaderSource(shader, 1, &strFileData, nil)

	gl.CompileShader(shader)

	var status int32;
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))
		fmt.Println(log)
		fmt.Errorf("failed to compile %v: %v", shaderText, log)
	}

	return shader
}
