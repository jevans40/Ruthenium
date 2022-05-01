package world

import (
	"math"
	"reflect"

	"github.com/jevans40/Ruthenium/linmath"
	"github.com/jevans40/Ruthenium/ruthutil"
)

type Coordinate [3]float64

func (c Coordinate) TranslateY(Y float64) Coordinate {
	c[1] += Y
	return c
}

func (c Coordinate) TranslateX(X float64) Coordinate {
	c[0] += X
	return c
}

func (c Coordinate) TranslateZ(Z float64) Coordinate {
	c[2] += Z
	return c
}

func (t Coordinate) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t Coordinate) IsComponent()          {}

func NewCoordinate() Coordinate {
	return [3]float64{0, 0, 0}
}

type Transform linmath.Mat3f[float64]

func (t Transform) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t Transform) IsComponent()          {}

func (t Transform) Translate(XCoord, YCoord float64) Transform {
	newMat := linmath.Identity[float64]()
	newMat[2] = XCoord
	newMat[5] = YCoord
	t = Transform((linmath.Mat3f[float64](t).MatMul((linmath.Mat3f[float64](newMat)))))
	return t
}

func (t Transform) Reflect() Transform {
	newMat := linmath.Identity[float64]()
	newMat[0] = -1
	t = Transform((linmath.Mat3f[float64](t).MatMul((linmath.Mat3f[float64](newMat)))))
	return t
}

func (t Transform) Scale(XScale, YScale float64) Transform {
	newMat := linmath.Identity[float64]()
	newMat[0] = XScale
	newMat[4] = YScale
	t = Transform((linmath.Mat3f[float64](t).MatMul((linmath.Mat3f[float64](newMat)))))
	return t
}

func (t Transform) Rotate(θ float64) Transform {
	newMat := linmath.Identity[float64]()
	newMat[0] = math.Cos(θ)
	newMat[1] = (-1) * math.Sin(θ)
	newMat[3] = math.Sin(θ)
	newMat[4] = math.Cos(θ)
	t = Transform((linmath.Mat3f[float64](t).MatMul((linmath.Mat3f[float64](newMat)))))
	return t
}

func (t Transform) Shear(XShear, YShear float64) Transform {
	newMat := linmath.Identity[float64]()
	newMat[1] = XShear
	newMat[3] = YShear
	t = Transform((linmath.Mat3f[float64](t).MatMul((linmath.Mat3f[float64](newMat)))))
	return t
}

func NewTransform() Transform {
	return Transform(linmath.Identity[float64]())
}

type Sprite struct {
	X     float32
	Y     float32
	H     float32
	W     float32
	M     float32
	Color [4]uint8
}

func NewTexturedSprite(xpos float32, ypos float32, texheight float32, texwidth float32, texmap float32, color ruthutil.Color) Sprite {
	return Sprite{xpos, ypos, texheight, texwidth, texmap, [4]uint8{color.Red, color.Green, color.Blue, color.Alpha}}
}

func NewUntexturedSprite(color ruthutil.Color) Sprite {
	return Sprite{0, 0, 0, 0, 0, [4]uint8{color.Red, color.Green, color.Blue, color.Alpha}}
}
func (t Sprite) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t Sprite) IsComponent()          {}
