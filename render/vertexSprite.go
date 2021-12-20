package render

var _ SpriteRendererSubscriber = (*VertexRenderable)(nil)

//A very simple renderable that the spriteRenderer can render, but has no real functionality
type VertexRenderable struct {
	verticeData []float32
	spriteNum   int32
	renderer    *SpriteRenderer
}

//Creates a new Sprite and automatically requests space from the renderer
func VertexSpriteFactory(thisRenderer *SpriteRenderer) *VertexRenderable {
	var vertdat []float32
	spritenum := 0
	sprite := VertexRenderable{vertdat, int32(spritenum), thisRenderer}
	sprite.verticeData, sprite.spriteNum = thisRenderer.SubscribeSprite(&sprite)
	return &sprite
}

func (thisSprite *VertexRenderable) SetVerticies(newdata []float32) {

	copy(thisSprite.verticeData[0:28], newdata)
}

func (thisSprite *VertexRenderable) RemoveSprite() {
	thisSprite.renderer.Unsubscribe(thisSprite, thisSprite.spriteNum)
}

func (thisSprite *VertexRenderable) UpdateRenderVert(returnSlice []float32, SpriteNum int32) {
	thisSprite.verticeData = returnSlice
	thisSprite.spriteNum = SpriteNum
}

func (thisSprite *VertexRenderable) RendererCallback() {

}
