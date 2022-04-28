package entity

import (
	"github.com/jevans40/Ruthenium/component"
	"github.com/jevans40/Ruthenium/world"
)

type EntityBuilder struct {
	data             []world.StorageWriteable
	numberOfEntities int
	callback         chan world.EntityCreationData
}

func NewEntityBuilder(Entity chan world.EntityCreationData) EntityBuilder {
	return EntityBuilder{callback: Entity}
}

func (e *EntityBuilder) AddComponent(NewComp world.StorageWriteable) *EntityBuilder {
	e.data = append(e.data, NewComp)
	return e
}

func (e *EntityBuilder) SetNumber(numberOfEntities int) *EntityBuilder {
	e.numberOfEntities = numberOfEntities
	return e
}

func (e *EntityBuilder) Build() chan component.EntityID {
	callback := make(chan component.EntityID, e.numberOfEntities)
	e.callback <- world.EntityCreationData{e.numberOfEntities, e.data, callback}
	return callback
}
