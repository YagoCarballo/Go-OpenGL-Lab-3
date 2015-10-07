package main

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	"./wrapper"
	"runtime"
	"unsafe"
)

var positionBufferObject uint32
var program uint32
var vao uint32

// This function is called before entering the main rendering loop.
// Use it for all you initialisation stuff
func initApp (glw *wrapper.Glw) {
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	vertexPositions := []float32{
		0.75, 0.75, 0.0, 1.0,
		0.75, -0.75, 0.0, 1.0,
		-0.75, -0.75, 0.0, 1.0,
	}

	gl.GenBuffers(1, &positionBufferObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexPositions), gl.Ptr(vertexPositions), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	tmpProgram, err := glw.LoadShader("shaders/basic.vert", "shaders/basic.frag")
	if err != nil {
		panic(err)
	} else {
		program = tmpProgram
	}
}

func display () {
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)

	gl.BindBuffer(gl.ARRAY_BUFFER, positionBufferObject)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 4, gl.FLOAT, false, 0, nil)

	gl.DrawArrays(gl.TRIANGLES, 0, 3)

	gl.DisableVertexAttribArray(0)
	gl.UseProgram(0)
}

func reshape (window *glfw.Window, width, height int) {
	gl.Viewport(0, 0, int32(width), int32(height))
}

func keyCallback (window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return;
	}

	fmt.Println("Key Pressed: ", key)

	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func main () {
	runtime.LockOSThread()

	glw := wrapper.NewGlw(1024, 768, "Hello Graphics World")
	if err := glw.OpenWindow(); err != nil {
		panic(err)
	}

	if err := gl.Init(); err != nil {
		panic(err)
	}

//	gl.DebugMessageCallback(debugProc, gl.Ptr(glw))

	glw.SetRenderer(display)
	glw.SetKeyCallback(keyCallback)
	glw.SetReshapeCallback(reshape)

	initApp(glw)

	glw.EventLoop()
}

func debugProc (source uint32,gltype uint32, id uint32, severity uint32, length int32, message string, userParam unsafe.Pointer) {
	fmt.Println("--> ", message)
}