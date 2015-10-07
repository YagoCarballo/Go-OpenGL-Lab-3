package wrapper

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"log"
	"fmt"
)

type Glw struct  {
	// Properties
	Width, Height int
	Title string

	// State
	fps int
	running bool
	Window *glfw.Window

	// Callbacks
	renderer func(glw *Glw)
	keyCallBack glfw.KeyCallback
	reshape glfw.FramebufferSizeCallback
}

// Constructor

func NewWrapper(width, height int, title string) *Glw {
	// Init GLFW
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}

	return &Glw{ width, height, title, 60, true, nil, nil, nil, nil }
}

// Public Functions

func (glw *Glw) CreateWindow () *glfw.Window {
	setOpenGlVersion()
	defer printOpenGlVersionInfo()

	// Creates the Window
	win, err := glfw.CreateWindow(glw.Width, glw.Height, glw.Title, nil, nil)
	if err != nil {
		panic(err)
	}

	// Initiates GL
	if err := gl.Init(); err != nil {
		panic(err)
	}

	// Sets this context as the current context
	win.MakeContextCurrent()

	// Sets the Clear Color (Background Color)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	// Enables Depth
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

//	win.SetInputMode(glfw.StickyKeysMode, 1)

	// Sets the Window to the Wrapper
	glw.SetWindow(win)
	return win
}

func (glw *Glw) StartLoop () {
	// If the Window is open keep looping
	for !glw.GetWindow().ShouldClose() {
		// Calls the Render Callback
		glw.renderer(glw)

		// Maintenance
		glfw.PollEvents()
		glw.GetWindow().SwapBuffers()
	}
}

func (glw *Glw) Terminate () {
	/* Clean up */
	glw.Window.Destroy()
	glfw.Terminate()
}

// Private Functions

func setOpenGlVersion() {
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)    // Necessary for OS X
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile) // Necessary for OS X
	glfw.WindowHint(glfw.OpenGLDebugContext, glfw.True)

	glfw.WindowHint(glfw.Resizable, glfw.True)
}

func printOpenGlVersionInfo() {
	fmt.Printf("%s\n", gl.GoStr(gl.GetString(gl.RENDERER)))
	fmt.Printf("%s\n", gl.GoStr(gl.GetString(gl.VERSION)))
}

// Getters & Setters

func (glw *Glw) SetFPS (fps int) {
	glw.fps = fps
}

func (glw *Glw) GetFPS () int {
	return glw.fps
}

func (glw *Glw) SetWindow (window *glfw.Window) {
	glw.Window = window
}

func (glw *Glw) GetWindow () *glfw.Window {
	return glw.Window
}

func (glw *Glw) SetRenderCallback (callback func(glw *Glw)) {
	glw.renderer = callback
}

func (glw *Glw) SetKeyCallBack (callback glfw.KeyCallback) {
	glw.keyCallBack = callback
	glw.Window.SetKeyCallback(callback)
}

func (glw *Glw) SetReshapeCallback (callback glfw.FramebufferSizeCallback) {
	glw.reshape = callback
	glw.Window.SetFramebufferSizeCallback(callback)
}