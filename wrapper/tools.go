package wrapper
import (
	"fmt"
	"strings"
	"io/ioutil"
	"github.com/go-gl/gl/all-core/gl"
	"github.com/kardianos/osext"
)

//
// Read File
// Reads a file and returns a String with it's contents.
// 						(string has a null pointer at the end (for OpenGL))
//
// @param path (string) the path to the file
//
// @return fileContents (string) the contents of the file as string
// @return error (error) the error (if any)
//
func ReadFile (path string) (string, error) {
	fileContents, err := ioutil.ReadFile(path)
	if err != nil {
		// If the file is missing, try to open it in the same folder as the current executable
		missing := strings.Contains(err.Error(), "no such file or directory")
		if missing {
			// Get the Folder of the current Executable
			dir, err := osext.ExecutableFolder()
			if err != nil {
				panic(err)
			}

			// Read the file and return content or error
			content, secondErr := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, path))
			return string(content) + "\x00", secondErr
		}
	}

	return string(fileContents) + "\x00", err
}

//
// File To String
// Reads a file and returns a String with it's contents.
// 						(string has a null pointer at the end (for OpenGL))
// 						(If an error happens it panics)
//
// @param path (string) the path to the file
//
// @return fileContents (string) the contents of the file as string
//
func FileToString (path string) string {
	content, err := ReadFile(path)
	if err != nil {
		panic(err)
	}

	return content
}

//
// Build Shader
// Creates and compiles a shader
//
// @param source (string) the path to the shader file
// @param shaderType (uint32) the shader type
//
// @return shader (uint32) the pointer to the shader
// @return error (error) the error (if any)
//
func BuildShader (source string, shaderType uint32) (uint32, error) {
	// Creates the Shader Object
	shader := gl.CreateShader(shaderType)

	// Reads the File
	fileContents := FileToString(source)

	// Converts the file contents into a valid C String
	csource := gl.Str(fileContents)

	// Loads the Shader's Source
	gl.ShaderSource(shader, 1, &csource, nil)

	// Compiles the Shader
	gl.CompileShader(shader)

	// Gets any errors that happened when building the Shader
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)

	// If there was an error, parse the C Error into a Go Error and return it
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	// Returns the shader if everything is OK
	return shader, nil
}

//
// Load Shader
// Load vertex and fragment shader and return the compiled program.
//
// @param vertexShaderSource (string) path to the vertex shader file
// @param fragmentShaderSource (string) path to the fragment shader file
//
// @return program (uint32) a pointer to the shader program
// @return error (error) the error (if any)
//
func LoadShader (vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	// Loads the Vertex shader file
	vertexShader, err := BuildShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	// Loads the fragment shader file
	fragmentShader, err := BuildShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	// Creates the Program
	program := gl.CreateProgram()

	// Attaches the Shaders to the program
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)

	// Links the program
	gl.LinkProgram(program)

	// Gets any error that happened when linking the program
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)

	// If there was any error, parse the C error and return it as a Go error
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	// Deletes the shaders
	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	// returns the program
	return program, nil
}
