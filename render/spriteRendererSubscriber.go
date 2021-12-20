package render

//This interface is for sprites to subscribe to the renderer.
//This will pass back an array where the sprite can update the state of its vertices.
//They must implement a way to deal with a change in renderer's vert array. Via Update()
type SpriteRendererSubscriber interface {
	UpdateRenderVert(returnSlice []float32, SpriteNum int32)

	//This function is called once per frame right before the scene is rendered
	RendererCallback()
}
