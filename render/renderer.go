package render

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/jevans40/Ruthenium/linmath"
)

//Renderer Is a interface for all render types, they have one job, given a window size, render shit
type Renderer interface {
	Render(Width, Height int32)
}

//SimpRenderer is the base renderer object
type SimpRenderer struct {
	vertexBufferObject  uint32
	vertexArrayObject   uint32
	elementBufferObject uint32
	numOfQuads          int32
	programObject       uint32
	vert                []float32
	elem                []uint32
	allocation          []bool
}

//Render Main render loop of the game
func Render(vbo, program uint32) {
	// 1st attribute buffer : vertices
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
	gl.EnableVertexAttribArray(0)

	//Send vertex information to the graphics card
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)

	gl.VertexAttribPointer(
		0,               // attribute 0. No particular reason for 0, but must match the layout in the shader.
		3,               // size
		gl.FLOAT,        // type
		false,           // normalized?
		0,               // stride
		gl.PtrOffset(0), // array buffer offset
	)

	// Draw the triangle !
	gl.DrawArrays(gl.TRIANGLES, 0, 3) // Starting from vertex 0; 3 vertices total -> 1 triangle
	gl.DisableVertexAttribArray(0)
}

//NewTestRenderer i dont know why i need a comment yet but here it be
func NewTestRenderer(programObject uint32) Renderer {
	vert1 := linmath.VerticeFactory(100, 100, 50, 0, 0, 0, 255, 255, 255, 0)
	vert2 := linmath.VerticeFactory(100, 1000, 50, 0, 0, 0, 255, 0, 255, 0)
	vert3 := linmath.VerticeFactory(1000, 100, 50, 0, 0, 200, 255, 255, 255, 0)
	vert4 := linmath.VerticeFactory(1000, 1000, 50, 0, 0, 200, 255, 0, 255, 0)

	triangle := make([]linmath.Vertice, 4)
	triangle[0] = vert1
	triangle[1] = vert2
	triangle[2] = vert3
	triangle[3] = vert4

	searialTriangle := []float32{}
	for _, v := range triangle {
		x := v.ToFloats()
		searialTriangle = append(searialTriangle, x[:]...)
	}

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	var vbo uint32

	// Generate 1 buffer, put the resulting identifier in vertexbuffer
	gl.GenBuffers(1, &vbo)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindVertexArray(0)

	newRenderer := SimpRenderer{vbo, vao, ebo, 1, programObject, searialTriangle, []uint32{0, 1, 2, 2, 3, 1}, []bool{true, false}}
	newRenderer.init()

	return &newRenderer
}

//Render - Rendering function for the simple renderer
func (thisRenderer *SimpRenderer) Render(width, height int32) {

	// 1st attribute buffer : vertices
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Viewport(0, 0, width, height)

	orthomat := linmath.NewOrthoMat4f(float32(height), 0, 0, float32(width), 1, 0)
	gl.UseProgram(thisRenderer.programObject)
	loc := gl.GetUniformLocation(thisRenderer.programObject, gl.Str("MVT"+"\x00"))
	mat := orthomat.ToFloats()
	gl.UniformMatrix4fv(loc, 1, false, &mat[0])
	thisRenderer.bind()
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0)) // Starting from vertex 0; 3 vertices total -> 1 triangle
	thisRenderer.unbind()

}

func (thisRenderer *SimpRenderer) init() {
	thisRenderer.bind()
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.0, 0.2, 0.7, 0.0)

	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)

	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 28, gl.PtrOffset(0))

	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 28, gl.PtrOffset(12))

	gl.VertexAttribPointer(2, 4, gl.UNSIGNED_BYTE, true, 28, gl.PtrOffset(20))

	gl.VertexAttribPointer(3, 1, gl.INT, false, 28, gl.PtrOffset(24))

	thisRenderer.unbind()
}

func (thisRenderer *SimpRenderer) bind() {
	gl.BindVertexArray(thisRenderer.vertexArrayObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, thisRenderer.vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(thisRenderer.vert)*4, gl.Ptr(thisRenderer.vert), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, thisRenderer.elementBufferObject)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(thisRenderer.elem)*4, gl.Ptr(thisRenderer.elem), gl.STATIC_DRAW)

}

func (thisRenderer *SimpRenderer) unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}
