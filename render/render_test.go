package render

import (
	"testing"
)

func TestSpriteCreation(t *testing.T) {
	/*
		thisRender := SpriteRendererFactory()
		sprite := SimpleSpriteFactory(&thisRender)
		sprite.Move(10, 10, 50)
		sprite.Resize(100, 100)
		sprite.Recolor(0, 255, 0, 255)

		sprite2 := SimpleSpriteFactory(&thisRender)
		sprite2.Move(200, 200, 50)
		sprite2.Resize(150, 200)
		sprite2.Recolor(255, 127, 0, 255)
	*/
	testAtlas := ImageAtlasFactory(4096, 20)
	testAtlas.AddImagesFromFolder("./../testres/image")
	//testImage := file.LoadImageFromFile("./../testres/image/God.png")
	//testAtlas.AddImageFromImage(testImage, "TestImage")
	testAtlas.Init()
	//file.SaveImageToFile("./../testres/image/AtlasTest.png", testAtlas.getAtlas(0))
	//messages := map[string](chan int){}

}
