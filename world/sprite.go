package world

import (
	"math"
	"reflect"

	"github.com/jevans40/Ruthenium/linmath"
	"github.com/jevans40/Ruthenium/ruthutil"
)

type Renderable struct {
	X float64
	Y float64
	Z float64

	H float64
	W float64

	TexX  float32
	TexY  float32
	TexH  float32
	TexW  float32
	TexM  float32
	Color [4]uint8

	verts [8]float64
}

func (t Renderable) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t Renderable) IsComponent()          {}

func NewRenderable() Renderable {
	return Renderable{verts: [8]float64{-0.5, -0.5, 0.5, -0.5, -0.5, 0.5, 0.5, 0.5}}
}

func (r Renderable) TranslateY(Y float64) Renderable {
	r.Y += Y
	return r
}

func (r Renderable) TranslateX(X float64) Renderable {
	r.X += X
	return r
}

func (r Renderable) TranslateZ(Z float64) Renderable {
	r.Z += Z
	return r
}

func (r Renderable) Translate(XCoord, YCoord float64) Renderable {
	newMat := linmath.Identity[float64]()
	newMat[2] = XCoord
	newMat[5] = YCoord
	for i := 0; i < len(r.verts)/2; i++ {
		x := newMat.VectorMul(r.verts[i*2], r.verts[i*2+1], r.Z)
		r.verts[i*2] = x[0]
		r.verts[i*2+1] = x[1]
		//X[2] is ignored
	}
	return r
}

func (r Renderable) Reflect() Renderable {
	newMat := linmath.Identity[float64]()
	newMat[0] = -1
	for i := 0; i < len(r.verts)/2; i++ {
		x := newMat.VectorMul(r.verts[i*2], r.verts[i*2+1], r.Z)
		r.verts[i*2] = x[0]
		r.verts[i*2+1] = x[1]
	}
	return r
}

func (r Renderable) Scale(XScale, YScale float64) Renderable {
	newMat := linmath.Identity[float64]()
	newMat[0] = XScale
	newMat[4] = YScale
	for i := 0; i < len(r.verts)/2; i++ {
		x := newMat.VectorMul(r.verts[i*2], r.verts[i*2+1], r.Z)
		r.verts[i*2] = x[0]
		r.verts[i*2+1] = x[1]
	}
	return r
}

func (r Renderable) Rotate(θ float64) Renderable {
	newMat := linmath.Identity[float64]()
	newMat[0] = math.Cos(θ)
	newMat[1] = (-1) * math.Sin(θ)
	newMat[3] = math.Sin(θ)
	newMat[4] = math.Cos(θ)
	for i := 0; i < len(r.verts)/2; i++ {
		x := newMat.VectorMul(r.verts[i*2], r.verts[i*2+1], r.Z)
		r.verts[i*2] = x[0]
		r.verts[i*2+1] = x[1]
	}
	return r
}

func (r Renderable) Shear(XShear, YShear float64) Renderable {
	newMat := linmath.Identity[float64]()
	newMat[1] = XShear
	newMat[3] = YShear
	for i := 0; i < len(r.verts)/2; i++ {
		x := newMat.VectorMul(r.verts[i*2], r.verts[i*2+1], r.Z)
		r.verts[i*2] = x[0]
		r.verts[i*2+1] = x[1]
	}
	return r
}

func (r Renderable) SetTexturedSprite(xpos float32, ypos float32, texheight float32, texwidth float32, texmap float32, color ruthutil.Color) Renderable {
	r.TexX = xpos
	r.TexY = ypos
	r.TexH = texheight
	r.TexW = texwidth
	r.TexM = texmap
	r.Color = [4]uint8{color.Red, color.Green, color.Blue, color.Alpha}
	return r
}

func (r Renderable) SetUntexturedSprite(color ruthutil.Color) Renderable {
	r.TexX = 0
	r.TexY = 0
	r.TexH = 0
	r.TexW = 0
	r.TexM = 0
	r.Color = [4]uint8{color.Red, color.Green, color.Blue, color.Alpha}
	return r
}

/////////////////////////////////////////////////////////////////
//////////////////* To replace *////////////////////////////////
///////////////////////////////////////////////////////////////
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
