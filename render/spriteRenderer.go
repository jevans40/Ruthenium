package render

import (
	"image"
	"image/draw"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/jevans40/Ruthenium/file"
	"github.com/jevans40/Ruthenium/linmath"
)

type SpriteRenderer struct {
	vertexBufferObject  uint32
	vertexArrayObject   uint32
	elementBufferObject uint32
	numOfSprites        int32
	maxSprites          int32
	programObject       uint32
	vert                []float32

	//The order of elements should be 0,1,2,1,2,3
	elem ElementBuffer

	allocation []bool

	uniformlocations map[string]map[uint32]int32

	//TODO: Make this a map
	subscribers []SpriteRendererSubscriber
}

func SpriteRendererFactory( /*SourceAtlas ImageAtlas,*/ ) SpriteRenderer {
	//size := SourceAtlas.imageSize
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	//This section is required, without it the game will crash.
	//I don't know why
	data := file.LoadImageFromFile("./res/image/God.png")
	rect := data.Bounds()
	rgba := image.NewRGBA(rect)
	draw.Draw(rgba, rect, data, rect.Min, draw.Src)
	//pixData := rgba.Pix
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(rect.Max.X-rect.Min.X), int32(rect.Max.Y-rect.Min.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))
	//First Generate a Vertex Array to bind the vertex buffer object and the element buffer object to
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	//Generate the vertex buffer object
	var vbo uint32
	gl.GenBuffers(1, &vbo)

	//Generate the element buffer object
	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindVertexArray(0)
	//Unbind the vertex array

	program, err := CreateDefaultProgram()
	if err != nil {
		panic(err)
	}
	var maxSprites = 1024
	var vert []float32
	var elem ElementBuffer
	var allocation []bool
	var subs []SpriteRendererSubscriber
	var uniform map[string]map[uint32]int32 = make(map[string]map[uint32]int32)
	//Make the sprite atlas

	newRenderer := SpriteRenderer{vbo, vao, ebo, 0, int32(maxSprites), program, vert, elem, allocation, uniform, subs}
	newRenderer.init()
	return newRenderer
}

func (thisRenderer *SpriteRenderer) Render(width, height int32) {
	//Setup uniforms only once
	if len(thisRenderer.uniformlocations) == 0 {
		thisRenderer.uniformlocations["texture1"+"\x00"] = make(map[uint32]int32)
		thisRenderer.uniformlocations["texture1"+"\x00"][thisRenderer.programObject] = gl.GetUniformLocation(thisRenderer.programObject, gl.Str("texture1"+"\x00"))
		thisRenderer.uniformlocations["MVT"+"\x00"] = make(map[uint32]int32)
		thisRenderer.uniformlocations["MVT"+"\x00"][thisRenderer.programObject] = gl.GetUniformLocation(thisRenderer.programObject, gl.Str("MVT"+"\x00"))
	}

	//Notify subscribers that the program is about to render.
	for _, v := range thisRenderer.subscribers {
		v.RendererCallback()
	}

	// 1st attribute buffer : vertices
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.Viewport(0, 0, width, height)

	orthomat := linmath.NewOrthoMat4f(float32(height), 0, 0, float32(width), 1, 0)
	gl.UseProgram(thisRenderer.programObject)
	gl.Uniform1i(thisRenderer.uniformlocations["texture1"+"\x00"][thisRenderer.programObject], 0)
	loc := thisRenderer.uniformlocations["MVT"+"\x00"][thisRenderer.programObject]
	mat := orthomat.ToFloats()
	gl.UniformMatrix4fv(loc, 1, false, &mat[0])
	thisRenderer.bind()
	gl.DrawElements(gl.TRIANGLES, 6*thisRenderer.numOfSprites, gl.UNSIGNED_INT, gl.PtrOffset(0)) // Starting from vertex 0; 3 vertices total -> 1 triangle
	thisRenderer.unbind()

}

func (thisRenderer *SpriteRenderer) init() {
	thisRenderer.allocation = make([]bool, thisRenderer.maxSprites)
	thisRenderer.vert = make([]float32, thisRenderer.maxSprites*28)
	thisRenderer.elem.setSize(int(thisRenderer.maxSprites))

	thisRenderer.bind()
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.5, 0.5, 0.5, 0.5)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

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

func (thisRenderer *SpriteRenderer) bind() {

	gl.BindVertexArray(thisRenderer.vertexArrayObject)
	gl.BindBuffer(gl.ARRAY_BUFFER, thisRenderer.vertexBufferObject)
	gl.BufferData(gl.ARRAY_BUFFER, len(thisRenderer.vert)*4, gl.Ptr(thisRenderer.vert), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, thisRenderer.elementBufferObject)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(thisRenderer.elem.GetArray())*4, gl.Ptr(thisRenderer.elem.GetArray()), gl.STATIC_DRAW)

}

func (thisRenderer *SpriteRenderer) unbind() {
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (thisRenderer *SpriteRenderer) SubscribeSprite(sprite SpriteRendererSubscriber) (returnSlice []float32, spriteNum int32) {
	spriteNum = int32(thisRenderer.allocate())
	thisRenderer.subscribers = append(thisRenderer.subscribers, sprite)
	returnSlice = thisRenderer.vert[spriteNum*28 : (spriteNum+1)*28]
	thisRenderer.numOfSprites = thisRenderer.numOfSprites + 1
	return
}

func (thisRenderer *SpriteRenderer) allocate() int32 {
	//Look for open space first
	for i := thisRenderer.numOfSprites; int(i) < len(thisRenderer.allocation); i++ {

		if !thisRenderer.allocation[i] {
			thisRenderer.allocation[i] = true
			return int32(i)
		}
	}

	//No room in current
	//Notify subscribers and generate a new array
	thisRenderer.maxSprites = thisRenderer.maxSprites * 2
	thisRenderer.elem.setSize(int(thisRenderer.maxSprites))
	thisRenderer.vert = make([]float32, thisRenderer.maxSprites*28)
	thisRenderer.allocation = make([]bool, thisRenderer.maxSprites)
	for i, v := range thisRenderer.subscribers {
		v.UpdateRenderVert(thisRenderer.vert[i*28:(i+1)*28], int32(i))
		thisRenderer.allocation[i] = true
	}
	return thisRenderer.allocate()
}

//Add check for 25% reduce array
func (thisRenderer *SpriteRenderer) Unsubscribe(sub SpriteRendererSubscriber, num int32) {
	//Zero data in the vertice
	newVert := make([]float32, 28)
	copy(thisRenderer.vert[num*28:(num+1)*28], newVert)

	//Deallocate
	thisRenderer.allocation[num] = false

	//Will be fixed with a proper map
	for i, v := range thisRenderer.subscribers {
		if v == sub {
			thisRenderer.subscribers[i] = thisRenderer.subscribers[len(thisRenderer.subscribers)-1]
			thisRenderer.subscribers = thisRenderer.subscribers[0 : len(thisRenderer.subscribers)-2]
			return
		}
	}
}
