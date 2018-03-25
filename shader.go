package main

import (
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type ShaderProgram struct {
	id uint32
}

func (sp *ShaderProgram) Delete() {
	if sp.id != 0 {
		log.Debugf("Deleting shader program 0x%08x", sp.id)
		gl.DeleteProgram(sp.id)
		sp.id = 0
	}
}

func (sp *ShaderProgram) GetUniformLocation(name string) int32 {
	str := gl.Str(name + "\000")

	return gl.GetUniformLocation(sp.id, str)
}

func (sp ShaderProgram) Use() {
	gl.UseProgram(sp.id)
}

func NewShaderProgram(vertexPath, fragmentPath string) (*ShaderProgram, error) {
	vertexShader, err := CompileShaderFile(vertexPath, gl.VERTEX_SHADER)
	if err != nil {
		return &ShaderProgram{}, err
	}
	defer vertexShader.Delete()

	fragmentShader, err := CompileShaderFile(fragmentPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return &ShaderProgram{}, err
	}
	defer fragmentShader.Delete()

	programId := gl.CreateProgram()

	vertexShader.AttachTo(programId)
	fragmentShader.AttachTo(programId)
	gl.LinkProgram(programId)

	var status int32
	gl.GetProgramiv(programId, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(programId, gl.INFO_LOG_LENGTH, &logLength)

		errLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(programId, logLength, nil, gl.Str(errLog))

		gl.DeleteProgram(programId)

		return &ShaderProgram{}, fmt.Errorf("failed to link program: %v", errLog)
	}

	program := &ShaderProgram{id: programId}
	runtime.SetFinalizer(program, func(sp *ShaderProgram){
		sp.Delete()
	})

	return program, nil
}




type Shader struct {
	id uint32
}

func (s Shader) AttachTo(spId uint32) {
	gl.AttachShader(spId, s.id)
}

func (s *Shader) Delete() {
	if s.id != 0 {
		log.Debugf("Deleting shader 0x%08x", s.id)
		gl.DeleteShader(s.id)
		s.id = 0
	}
}

func CompileShaderFile(path string, shaderType uint32) (*Shader, error) {
	shaderSource, err := ioutil.ReadFile(path)
	if err != nil {
		return &Shader{}, err
	}

	return CompileShader(string(shaderSource), shaderType)
}

func CompileShader(source string, shaderType uint32) (*Shader, error) {
	shaderId := gl.CreateShader(shaderType)

	csource, csourceFree := gl.Strs(source)
	defer csourceFree()
	ln := int32(len(source))
	gl.ShaderSource(shaderId, 1, csource, &ln)
	gl.CompileShader(shaderId)

	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)

		errLog := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(errLog))

		return &Shader{}, fmt.Errorf("failed to compile shader %v: %v", source, errLog)
	}

	shader := &Shader{id: shaderId}
	runtime.SetFinalizer(shader, func(s *Shader){
		s.Delete()
	})

	return shader, nil
}
