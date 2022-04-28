package world

import (
	"reflect"

	"github.com/jevans40/Ruthenium/ruthutil"
)

type Transform struct {
	X float32
	Y float32
	Z float32

	W float32
	H float32
}

func (t Transform) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t Transform) IsComponent()          {}
func NewTransform(xpos float32, ypos float32, zpos float32, width float32, height float32) Transform {
	return Transform{xpos, ypos, zpos, width, height}
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
