package render

//TODO:: Concider documentation
type Sprite interface {

	//Mutators
	Resize(x float32, y float32)
	Move(x, y, z float32)
	Recolor(r uint8, g uint8, b uint8, a uint8)
	SetTexPos(x, y float32)
	SetTexSize(x, y float32)
	SetMap(newMap uint32)

	//Accessors
	GetSize() (x, y float32)
	GetPos() (x, y, z float32)
	GetColor() (r, g, b, a uint8)
	GetTexPos() (x, y float32)
	GetTexSize() (x, y float32)
	GetMap() (texMap uint32)
	//Resize
	//
}
