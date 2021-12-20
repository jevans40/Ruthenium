package component

import "reflect"

type Component interface {
	IsComponent()
	GetType() reflect.Type
}

var _ Component = &BaseComponent{}

type BaseComponent struct {
}

func (BaseComponent) IsComponent()            {}
func (b BaseComponent) GetType() reflect.Type { return reflect.TypeOf((*BaseComponent)(nil)).Elem() }

func ReflectType[T any]() reflect.Type { return reflect.TypeOf((*T)(nil)).Elem() }
