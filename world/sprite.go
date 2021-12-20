package world

import "reflect"

type Transform struct {
	X float32
	Y float32
	Z float32

	W float32
	H float32
}

func (t Transform) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t Transform) IsComponent()          {}

type Sprite struct {
	X     float32
	Y     float32
	H     float32
	W     float32
	M     float32
	Color [4]uint8
}

func (t Sprite) GetType() reflect.Type { return reflect.TypeOf(t) }
func (t Sprite) IsComponent()          {}
