package main

import (
	"os"
	"runtime"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"

	// //"math"
	// mgl "github.com/go-gl/mathgl/mgl64"

	"math"

	"github.com/op/go-logging"
)

var (
	log     = logging.MustGetLogger("frax")
	nullptr = unsafe.Pointer(uintptr(0))
)

func init() {
	logFormat := logging.MustStringFormatter("%{time:15:04:05.000} (%{module}) %{color}%{level:.4s}%{color:reset} %{message}")
	logging.SetFormatter(logFormat)

	backend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(backend)

	runtime.LockOSThread()
}

func ptr(i int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(i))
}

func main() {
	log.Debug("Initializing GLFW")
	err := glfw.Init()
	if err != nil {
		log.Critical(err)
	}
	defer glfw.Terminate()

	log.Debug("Creating main window")
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	// glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.Resizable, glfw.True)

	window, err := glfw.CreateWindow(600, 600, "FraX", nil, nil)
	if err != nil {
		log.Critical(err)
	}
	window.MakeContextCurrent()

	log.Debug("Initializing OpenGL")
	err = gl.Init()
	if err != nil {
		log.Critical(err)
	}

	//output glfw.SwapInterval(0)
	glfw.SwapInterval(1)

	width, height := window.GetFramebufferSize()
	gl.Viewport(0, 0, int32(width), int32(height))

	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Info("OpenGL version :", version)

	log.Debug("Loading shaders...")
	shaderProg, err := NewShaderProgram("vertex.glsl", "fragment.glsl")
	if err != nil {
		log.Fatal(err)
	}
	defer shaderProg.Delete()
	log.Debug("Shaders loaded")

	var max_it int32 = 10

	window.SetKeyCallback(func(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if key == glfw.KeyEscape && action == glfw.Press {
			log.Notice("User pressed Esc key")
			window.SetShouldClose(true)
		}

		if key == glfw.KeyUp && (action == glfw.Press || action == glfw.Repeat) {
			max_it++
			log.Infof("max_it = %d", max_it)
		}

		if key == glfw.KeyDown && (action == glfw.Press || action == glfw.Repeat) {
			max_it--
			log.Infof("max_it = %d", max_it)
		}
	})

	scale := 2.0
	xAbsOff := -10.0
	window.SetScrollCallback(func(window *glfw.Window, xoff, yoff float64) {
		xAbsOff += yoff
		scale = math.Pow(2.0, -xAbsOff/10.0)
	})

	gl.Viewport(0, 0, 600, 600)
	s := 600
	window.SetSizeCallback(func(window *glfw.Window, width, height int) {
		s = width
		if width < height {
			s = height
		}

		gl.Viewport(int32((width-s)/2), int32((height-s)/2), int32(s), int32(s))
	})

	vertices := []float64{
		// Positions
		-1.0, 01.0, 0.0, // Top Left
		01.0, 01.0, 0.0, // Top Right
		01.0, -1.0, 0.0, // Bottom Right
		-1.0, -1.0, 0.0, // Bottom Left
	}

	indices := []uint32{
		0, 1, 2,
		0, 2, 3,
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	defer gl.DeleteVertexArrays(1, &vao)
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &vbo)
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	defer gl.DeleteBuffers(1, &ebo)
	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*8, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	gl.VertexAttribPointer(0, 3, gl.DOUBLE, false, 3*8, nullptr)
	gl.EnableVertexAttribArray(0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

	log.Debug("Starting main loop")
	dragL := false
	dragR := false
	dragStartX, dragStartY := 0.0, 0.0
	centerX, centerY := 0.0, 0.0
	centerBackX, centerBackY := 0.0, 0.0
	cRe, cIm := 0.0, 0.0
	cBackRe, cBackIm := 0.0, 0.0
	for !window.ShouldClose() {
		if window.GetMouseButton(glfw.MouseButtonLeft) == glfw.Press {
			if !dragL {
				dragL = true
				dragStartX, dragStartY = window.GetCursorPos()
				centerBackX = centerX
				centerBackY = centerY
			}
		} else if window.GetMouseButton(glfw.MouseButtonRight) == glfw.Press {
			if !dragR {
				dragR = true
				dragStartX, dragStartY = window.GetCursorPos()
				cBackRe = cRe
				cBackIm = cIm
			}
		} else {
			dragL = false
			dragR = false
		}

		if dragL {
			cursorX, cursorY := window.GetCursorPos()
			centerX = centerBackX - (cursorX-dragStartX)*scale
			centerY = centerBackY + (cursorY-dragStartY)*scale
		}

		if dragR {
			cursorX, cursorY := window.GetCursorPos()
			cRe = cBackRe - (cursorX-dragStartX)/400*scale
			cIm = cBackIm + (cursorY-dragStartY)/400*scale
		}

		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT) // | gl.DEPTH_BUFFER_BIT)

		shaderProg.Use()

		gl.Uniform2f(shaderProg.GetUniformLocation("center"), float32(centerX/float64(s)*2.0), float32(centerY/float64(s)*2.0))
		gl.Uniform1f(shaderProg.GetUniformLocation("scale"), float32(scale))
		// gl.Uniform2f(shaderProg.GetUniformLocation("c"), float32(0.0), float32(1.0))
		// gl.Uniform2f(shaderProg.GetUniformLocation("c"), float32(-0.6), float32(0.6))
		gl.Uniform2f(shaderProg.GetUniformLocation("c"), float32(cRe), float32(cIm))
		// gl.Uniform2f(shaderProg.GetUniformLocation("c"), float32(0.7885*math.Cos(glfw.GetTime()/5)), float32(0.7885*math.Sin(glfw.GetTime()/5)))

		gl.Uniform1i(shaderProg.GetUniformLocation("max_it"), max_it)

		gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nullptr)
		gl.BindVertexArray(0)

		window.SwapBuffers()

		glfw.PollEvents()
	}
	log.Debug("Window should close")

	log.Notice("Exiting")
}
