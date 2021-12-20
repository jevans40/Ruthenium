package linmath

type PSPoint interface {
	X() int32
	Y() int32
	Setx(x int32)
	Sety(y int32)
}

//This is just a simple rectangle, used for within the PS engine
type PSRectangle interface {
	//Returns a size touple of the rectangle
	GetSize() PSPoint
	//Returns the top right corner of the rectangle
	GetPoint() PSPoint

	//Set the size of the rectangle
	SetSize(size PSPoint)

	//Set the top right corner coordinates of the rectangle
	SetPoint(point PSPoint)
}

type rectangle [4]int32
type point [2]int32

func NewPSPoint(x, y int32) PSPoint {
	return &point{x, y}
}

func (i *point) X() int32 {
	return i[0]
}
func (i *point) Y() int32 {
	return i[1]
}

func (i *point) Setx(x int32) {
	i[0] = x
}

func (i *point) Sety(y int32) {
	i[0] = y
}

func NewPSRectangle(sx, sy, x, y int32) PSRectangle {
	return &rectangle{sx, sy, x, y}
}

func (i *rectangle) GetSize() PSPoint {
	return &point{i[0], i[1]}
}

func (i *rectangle) GetPoint() PSPoint {
	return &point{i[2], i[3]}
}

func (i *rectangle) SetSize(size PSPoint) {
	i[0] = size.X()
	i[1] = size.Y()
}

func (i *rectangle) SetPoint(point PSPoint) {
	i[2] = point.X()
	i[3] = point.Y()
}
