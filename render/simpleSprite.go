package render

import "github.com/jevans40/Ruthenium/linmath"

var _ SpriteRendererSubscriber = (*simpleSprite)(nil)
var _ Sprite = (*simpleSprite)(nil)

type simpleSprite struct {
	verticeData     []float32
	size            [2]float32
	position        [3]float32
	texturePosition [2]float32
	textureSize     [2]float32
	textureMap      uint32
	color           [4]uint8
	spriteNum       int32
	SpriteRender    *SpriteRenderer
}

//Creates a new Sprite and automatically requests space from the renderer
func SimpleSpriteFactory(thisRenderer *SpriteRenderer) Sprite {
	var vertdat []float32
	size := [2]float32{0, 0}
	pos := [3]float32{0, 0, 0}
	texPos := [2]float32{0, 0}
	texSize := [2]float32{0, 0}
	texMap := 0
	col := [4]uint8{255, 0, 0, 0}
	spritenum := 0

	sprite := simpleSprite{vertdat, size, pos, texPos, texSize, uint32(texMap), col, int32(spritenum), thisRenderer}
	sprite.verticeData, sprite.spriteNum = thisRenderer.SubscribeSprite(&sprite)
	sprite.calculateVerticies()
	return &sprite
}

func (thisSprite *simpleSprite) calculateVerticies() {

	vert := make([]linmath.Vertice, 4)
	for i, _ := range vert {
		vert[i] = linmath.EmptyVertice()
	}

	for i := 0; i < 4; i++ {
		vert[i].SetColor(thisSprite.color)
		vert[i].SetMap(thisSprite.textureMap)
		vert[i].SetTexX(thisSprite.texturePosition[0] + thisSprite.textureSize[0]*float32((i)%2))
		vert[i].SetTexY(thisSprite.texturePosition[1] + thisSprite.textureSize[1]*float32(int32((i)/2)%2))
		vert[i].SetX(thisSprite.position[0] + thisSprite.size[0]*float32((i)%2))
		vert[i].SetY(thisSprite.position[1] + thisSprite.size[1]*float32(int32((i)/2)%2))
		vert[i].SetZ(thisSprite.position[2])
	}

	for i, v := range vert {
		floats := v.ToFloats()
		copy(thisSprite.verticeData[7*i:7*i+7], floats[:])
	}
}

func (thisSprite *simpleSprite) RemoveSprite() {
	thisSprite.SpriteRender.Unsubscribe(thisSprite, thisSprite.spriteNum)
}

func (thisSprite *simpleSprite) Resize(x float32, y float32) {
	thisSprite.size = [2]float32{x, y}
}

func (thisSprite *simpleSprite) Move(x float32, y float32, z float32) {
	thisSprite.position = [3]float32{x, y, z}
}

func (thisSprite *simpleSprite) Recolor(r uint8, g uint8, b uint8, a uint8) {
	thisSprite.color = [4]uint8{r, g, b, a}
}

func (thisSprite *simpleSprite) SetTexPos(x float32, y float32) {
	thisSprite.texturePosition = [2]float32{x, y}
}

func (thisSprite *simpleSprite) SetTexSize(x float32, y float32) {
	thisSprite.textureSize = [2]float32{x, y}
}

func (thisSprite *simpleSprite) SetMap(newMap uint32) {
	thisSprite.textureMap = newMap
}

func (thisSprite *simpleSprite) GetSize() (x, y float32) {
	x, y = thisSprite.size[0], thisSprite.size[1]
	return
}

func (thisSprite *simpleSprite) GetPos() (x, y, z float32) {
	x, y, z = thisSprite.position[0], thisSprite.position[1], thisSprite.position[2]
	return
}

func (thisSprite *simpleSprite) GetColor() (r, g, b, a uint8) {
	r, g, b, a = thisSprite.color[0], thisSprite.color[1], thisSprite.color[2], thisSprite.color[3]
	return
}

func (thisSprite *simpleSprite) GetTexPos() (x, y float32) {
	x, y = thisSprite.texturePosition[0], thisSprite.texturePosition[1]
	return
}

func (thisSprite *simpleSprite) GetTexSize() (x, y float32) {
	x, y = thisSprite.textureSize[0], thisSprite.textureSize[1]
	return
}

func (thisSprite *simpleSprite) GetMap() (texMap uint32) {
	texMap = thisSprite.textureMap
	return
}

func (thisSprite *simpleSprite) UpdateRenderVert(returnSlice []float32, SpriteNum int32) {
	thisSprite.verticeData = returnSlice
	thisSprite.spriteNum = SpriteNum
}

func (thisSprite *simpleSprite) RendererCallback() {
	thisSprite.calculateVerticies()
}
